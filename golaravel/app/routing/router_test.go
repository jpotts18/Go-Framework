package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golaravel/app/http/request"
	"golaravel/app/http/response"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Fatal("Expected router to be created")
	}
	if r.routes == nil {
		t.Error("Expected routes slice to be initialized")
	}
	if r.namedRoutes == nil {
		t.Error("Expected namedRoutes map to be initialized")
	}
}

func TestRouterGet(t *testing.T) {
	r := NewRouter()

	route := r.Get("/test", func(req *request.Request, res *response.Response) error {
		return res.String("OK")
	})

	if route.Method != "GET" {
		t.Errorf("Expected method to be GET, got: %s", route.Method)
	}
	if route.Pattern != "/test" {
		t.Errorf("Expected pattern to be /test, got: %s", route.Pattern)
	}
}

func TestRouterPost(t *testing.T) {
	r := NewRouter()

	route := r.Post("/test", func(req *request.Request, res *response.Response) error {
		return nil
	})

	if route.Method != "POST" {
		t.Errorf("Expected method to be POST, got: %s", route.Method)
	}
}

func TestRouterPut(t *testing.T) {
	r := NewRouter()

	route := r.Put("/test", func(req *request.Request, res *response.Response) error {
		return nil
	})

	if route.Method != "PUT" {
		t.Errorf("Expected method to be PUT, got: %s", route.Method)
	}
}

func TestRouterDelete(t *testing.T) {
	r := NewRouter()

	route := r.Delete("/test", func(req *request.Request, res *response.Response) error {
		return nil
	})

	if route.Method != "DELETE" {
		t.Errorf("Expected method to be DELETE, got: %s", route.Method)
	}
}

func TestRouterRouteParameters(t *testing.T) {
	r := NewRouter()

	var capturedID string
	r.Get("/users/{id}", func(req *request.Request, res *response.Response) error {
		capturedID = req.Param("id")
		return res.String("OK")
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if capturedID != "123" {
		t.Errorf("Expected param id to be '123', got: %s", capturedID)
	}
}

func TestRouterMultipleParameters(t *testing.T) {
	r := NewRouter()

	var userID, postID string
	r.Get("/users/{userId}/posts/{postId}", func(req *request.Request, res *response.Response) error {
		userID = req.Param("userId")
		postID = req.Param("postId")
		return res.String("OK")
	})

	req := httptest.NewRequest("GET", "/users/42/posts/99", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if userID != "42" {
		t.Errorf("Expected userId to be '42', got: %s", userID)
	}
	if postID != "99" {
		t.Errorf("Expected postId to be '99', got: %s", postID)
	}
}

func TestRouterNamedRoutes(t *testing.T) {
	r := NewRouter()

	r.Get("/users/{id}", func(req *request.Request, res *response.Response) error {
		return nil
	}).Name("user.show")

	route := r.Route("user.show")
	if route == nil {
		t.Fatal("Expected named route to be found")
	}
	if route.Pattern != "/users/{id}" {
		t.Errorf("Expected pattern to be /users/{id}, got: %s", route.Pattern)
	}
}

func TestRouterURL(t *testing.T) {
	r := NewRouter()

	r.Get("/users/{id}", func(req *request.Request, res *response.Response) error {
		return nil
	}).Name("user.show")

	url := r.URL("user.show", map[string]string{"id": "42"})
	if url != "/users/42" {
		t.Errorf("Expected URL to be /users/42, got: %s", url)
	}
}

func TestRouterGroup(t *testing.T) {
	r := NewRouter()

	r.Group("/api", func(group *RouteGroup) {
		group.Get("/users", func(req *request.Request, res *response.Response) error {
			return res.String("API Users")
		})
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", w.Code)
	}
}

func TestRouterNestedGroups(t *testing.T) {
	r := NewRouter()

	r.Group("/api", func(api *RouteGroup) {
		api.Group("/v1", func(v1 *RouteGroup) {
			v1.Get("/users", func(req *request.Request, res *response.Response) error {
				return res.String("V1 Users")
			})
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", w.Code)
	}
	if w.Body.String() != "V1 Users" {
		t.Errorf("Expected body 'V1 Users', got: %s", w.Body.String())
	}
}

func TestRouterMiddleware(t *testing.T) {
	r := NewRouter()

	middlewareCalled := false
	testMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			middlewareCalled = true
			return next(req, res)
		}
	}

	r.Use(testMiddleware)
	r.Get("/test", func(req *request.Request, res *response.Response) error {
		return res.String("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Expected middleware to be called")
	}
}

func TestRouterRouteMiddleware(t *testing.T) {
	r := NewRouter()

	routeMiddlewareCalled := false
	routeMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			routeMiddlewareCalled = true
			return next(req, res)
		}
	}

	r.Get("/test", func(req *request.Request, res *response.Response) error {
		return res.String("OK")
	}, routeMiddleware)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !routeMiddlewareCalled {
		t.Error("Expected route middleware to be called")
	}
}

func TestRouterGroupMiddleware(t *testing.T) {
	r := NewRouter()

	groupMiddlewareCalled := false
	groupMiddleware := func(next HandlerFunc) HandlerFunc {
		return func(req *request.Request, res *response.Response) error {
			groupMiddlewareCalled = true
			return next(req, res)
		}
	}

	r.Group("/api", func(group *RouteGroup) {
		group.Use(groupMiddleware)
		group.Get("/test", func(req *request.Request, res *response.Response) error {
			return res.String("OK")
		})
	})

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !groupMiddlewareCalled {
		t.Error("Expected group middleware to be called")
	}
}

func TestRouterNotFound(t *testing.T) {
	r := NewRouter()

	r.Get("/existing", func(req *request.Request, res *response.Response) error {
		return res.String("OK")
	})

	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got: %d", w.Code)
	}
}

func TestRouterCustomNotFound(t *testing.T) {
	r := NewRouter()

	r.NotFound(func(req *request.Request, res *response.Response) error {
		return res.Status(404).String("Custom Not Found")
	})

	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got: %d", w.Code)
	}
	if w.Body.String() != "Custom Not Found" {
		t.Errorf("Expected 'Custom Not Found', got: %s", w.Body.String())
	}
}

func TestRouterRoutes(t *testing.T) {
	r := NewRouter()

	r.Get("/route1", func(req *request.Request, res *response.Response) error { return nil })
	r.Post("/route2", func(req *request.Request, res *response.Response) error { return nil })

	routes := r.Routes()
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got: %d", len(routes))
	}
}

func TestRouterAny(t *testing.T) {
	r := NewRouter()

	r.Any("/any", func(req *request.Request, res *response.Response) error {
		return res.String("Any method")
	})

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/any", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for %s, got: %d", method, w.Code)
		}
	}
}

func TestRouterMatch(t *testing.T) {
	r := NewRouter()

	r.Match([]string{"GET", "POST"}, "/match", func(req *request.Request, res *response.Response) error {
		return res.String("Matched")
	})

	getReq := httptest.NewRequest("GET", "/match", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected GET to match, got status: %d", getW.Code)
	}

	postReq := httptest.NewRequest("POST", "/match", nil)
	postW := httptest.NewRecorder()
	r.ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Errorf("Expected POST to match, got status: %d", postW.Code)
	}

	deleteReq := httptest.NewRequest("DELETE", "/match", nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusNotFound {
		t.Errorf("Expected DELETE to not match, got status: %d", deleteW.Code)
	}
}
