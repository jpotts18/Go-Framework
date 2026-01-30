package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golaravel/app/http/request"
	"golaravel/app/http/response"
	"golaravel/app/routing"
)

func Logger() routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			start := time.Now()
			err := next(req, res)
			duration := time.Since(start)
			log.Printf("[%s] %s %s - %v", req.Method(), req.Path(), req.IP(), duration)
			return err
		}
	}
}

func Recovery() routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic recovered: %v", r)
					err = res.ServerError("Internal Server Error")
				}
			}()
			return next(req, res)
		}
	}
}

func CORS(allowedOrigins ...string) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			origin := req.Header("Origin")

			allowed := false
			if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*") {
				allowed = true
				res.Header("Access-Control-Allow-Origin", "*")
			} else {
				for _, o := range allowedOrigins {
					if o == origin {
						allowed = true
						res.Header("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if allowed {
				res.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				res.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, X-Requested-With")
				res.Header("Access-Control-Allow-Credentials", "true")
				res.Header("Access-Control-Max-Age", "86400")
			}

			if req.Method() == "OPTIONS" {
				return res.NoContent()
			}

			return next(req, res)
		}
	}
}

func JSON() routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			res.Header("Content-Type", "application/json")
			return next(req, res)
		}
	}
}

func Auth(validateToken func(string) bool) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			token := req.BearerToken()
			if token == "" {
				return res.Unauthorized("Missing authentication token")
			}

			if !validateToken(token) {
				return res.Unauthorized("Invalid authentication token")
			}

			return next(req, res)
		}
	}
}

func BasicAuth(users map[string]string) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			username, password, ok := req.Raw().BasicAuth()
			if !ok {
				res.Header("WWW-Authenticate", `Basic realm="Restricted"`)
				return res.Unauthorized("Authentication required")
			}

			expectedPassword, exists := users[username]
			if !exists || password != expectedPassword {
				res.Header("WWW-Authenticate", `Basic realm="Restricted"`)
				return res.Unauthorized("Invalid credentials")
			}

			return next(req, res)
		}
	}
}

func RateLimiter(maxRequests int, window time.Duration) routing.MiddlewareFunc {
	requests := make(map[string][]time.Time)

	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			ip := req.IP()
			now := time.Now()

			times, exists := requests[ip]
			if !exists {
				times = make([]time.Time, 0)
			}

			validTimes := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < window {
					validTimes = append(validTimes, t)
				}
			}

			if len(validTimes) >= maxRequests {
				res.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
				return res.Status(429).JSON(map[string]string{
					"error": "Too many requests",
				})
			}

			validTimes = append(validTimes, now)
			requests[ip] = validTimes

			return next(req, res)
		}
	}
}

func Timeout(duration time.Duration) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			done := make(chan error, 1)
			go func() {
				done <- next(req, res)
			}()

			select {
			case err := <-done:
				return err
			case <-time.After(duration):
				return res.Status(504).JSON(map[string]string{
					"error": "Request timeout",
				})
			}
		}
	}
}

func TrimStrings() routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			return next(req, res)
		}
	}
}

func SecureHeaders() routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			res.Header("X-Content-Type-Options", "nosniff")
			res.Header("X-Frame-Options", "SAMEORIGIN")
			res.Header("X-XSS-Protection", "1; mode=block")
			res.Header("Referrer-Policy", "strict-origin-when-cross-origin")
			return next(req, res)
		}
	}
}

func ContentType(contentType string) routing.MiddlewareFunc {
	return func(next routing.HandlerFunc) routing.HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			reqContentType := req.Header("Content-Type")
			if req.Method() != "GET" && req.Method() != "HEAD" && req.Method() != "OPTIONS" {
				if !strings.Contains(reqContentType, contentType) {
					return res.Status(415).JSON(map[string]string{
						"error": "Unsupported Media Type",
					})
				}
			}
			return next(req, res)
		}
	}
}
