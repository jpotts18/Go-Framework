package routing

import (
        "net/http"
        "regexp"
        "strings"

        "golaravel/app/http/request"
        "golaravel/app/http/response"
)

type HandlerFunc func(*request.Request, *response.Response) error

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Route struct {
        Method      string
        Pattern     string
        Handler     HandlerFunc
        Middleware  []MiddlewareFunc
        name        string
        regex       *regexp.Regexp
        paramNames  []string
}

type RouteGroup struct {
        router     *Router
        prefix     string
        middleware []MiddlewareFunc
}

type Router struct {
        routes          []*Route
        middleware      []MiddlewareFunc
        notFoundHandler HandlerFunc
        namedRoutes     map[string]*Route
}

func NewRouter() *Router {
        return &Router{
                routes:      make([]*Route, 0),
                middleware:  make([]MiddlewareFunc, 0),
                namedRoutes: make(map[string]*Route),
        }
}

func (r *Router) compilePattern(pattern string) (*regexp.Regexp, []string) {
        paramNames := make([]string, 0)

        regexPattern := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllStringFunc(pattern, func(match string) string {
                paramName := strings.Trim(match, "{}")

                if strings.HasSuffix(paramName, "?") {
                        paramName = strings.TrimSuffix(paramName, "?")
                        paramNames = append(paramNames, paramName)
                        return `([^/]*)?`
                }

                paramNames = append(paramNames, paramName)
                return `([^/]+)`
        })

        regexPattern = "^" + regexPattern + "$"
        regex := regexp.MustCompile(regexPattern)

        return regex, paramNames
}

func (r *Router) addRoute(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        regex, paramNames := r.compilePattern(pattern)

        route := &Route{
                Method:     method,
                Pattern:    pattern,
                Handler:    handler,
                Middleware: middleware,
                regex:      regex,
                paramNames: paramNames,
        }

        r.routes = append(r.routes, route)
        return route
}

func (r *Router) Get(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("GET", pattern, handler, middleware...)
}

func (r *Router) Post(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("POST", pattern, handler, middleware...)
}

func (r *Router) Put(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("PUT", pattern, handler, middleware...)
}

func (r *Router) Patch(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("PATCH", pattern, handler, middleware...)
}

func (r *Router) Delete(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("DELETE", pattern, handler, middleware...)
}

func (r *Router) Options(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return r.addRoute("OPTIONS", pattern, handler, middleware...)
}

func (r *Router) Any(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
        methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
        for _, method := range methods {
                r.addRoute(method, pattern, handler, middleware...)
        }
}

func (r *Router) Match(methods []string, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
        for _, method := range methods {
                r.addRoute(strings.ToUpper(method), pattern, handler, middleware...)
        }
}

func (route *Route) Name(name string) *Route {
        route.name = name
        return route
}

func (route *Route) GetName() string {
        return route.name
}

func (r *Router) Use(middleware ...MiddlewareFunc) {
        r.middleware = append(r.middleware, middleware...)
}

func (r *Router) Group(prefix string, callback func(*RouteGroup)) *RouteGroup {
        group := &RouteGroup{
                router:     r,
                prefix:     prefix,
                middleware: make([]MiddlewareFunc, 0),
        }
        callback(group)
        return group
}

func (g *RouteGroup) Use(middleware ...MiddlewareFunc) *RouteGroup {
        g.middleware = append(g.middleware, middleware...)
        return g
}

func (g *RouteGroup) addRoute(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        fullPattern := g.prefix + pattern
        allMiddleware := append(g.middleware, middleware...)
        return g.router.addRoute(method, fullPattern, handler, allMiddleware...)
}

func (g *RouteGroup) Get(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return g.addRoute("GET", pattern, handler, middleware...)
}

func (g *RouteGroup) Post(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return g.addRoute("POST", pattern, handler, middleware...)
}

func (g *RouteGroup) Put(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return g.addRoute("PUT", pattern, handler, middleware...)
}

func (g *RouteGroup) Patch(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return g.addRoute("PATCH", pattern, handler, middleware...)
}

func (g *RouteGroup) Delete(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
        return g.addRoute("DELETE", pattern, handler, middleware...)
}

func (g *RouteGroup) Group(prefix string, callback func(*RouteGroup)) *RouteGroup {
        newGroup := &RouteGroup{
                router:     g.router,
                prefix:     g.prefix + prefix,
                middleware: append([]MiddlewareFunc{}, g.middleware...),
        }
        callback(newGroup)
        return newGroup
}

func (r *Router) NotFound(handler HandlerFunc) {
        r.notFoundHandler = handler
}

func (r *Router) matchRoute(method, path string) (*Route, map[string]string) {
        for _, route := range r.routes {
                if route.Method != method {
                        continue
                }

                matches := route.regex.FindStringSubmatch(path)
                if matches == nil {
                        continue
                }

                params := make(map[string]string)
                for i, name := range route.paramNames {
                        if i+1 < len(matches) {
                                params[name] = matches[i+1]
                        }
                }

                if route.name != "" {
                        r.namedRoutes[route.name] = route
                }

                return route, params
        }

        return nil, nil
}

func (r *Router) applyMiddleware(handler HandlerFunc, middleware []MiddlewareFunc) HandlerFunc {
        for i := len(middleware) - 1; i >= 0; i-- {
                handler = middleware[i](handler)
        }
        return handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
        httpReq := request.New(req)
        httpRes := response.New(w)

        route, params := r.matchRoute(req.Method, req.URL.Path)

        if route == nil {
                if r.notFoundHandler != nil {
                        r.notFoundHandler(httpReq, httpRes)
                        return
                }
                httpRes.NotFound("Route not found")
                return
        }

        for key, value := range params {
                httpReq.SetParam(key, value)
        }

        allMiddleware := append(r.middleware, route.Middleware...)
        handler := r.applyMiddleware(route.Handler, allMiddleware)

        if err := handler(httpReq, httpRes); err != nil {
                httpRes.ServerError(err.Error())
        }
}

func (r *Router) Route(name string) *Route {
        return r.namedRoutes[name]
}

func (r *Router) URL(name string, params map[string]string) string {
        route := r.namedRoutes[name]
        if route == nil {
                return ""
        }

        url := route.Pattern
        for key, value := range params {
                url = strings.Replace(url, "{"+key+"}", value, 1)
                url = strings.Replace(url, "{"+key+"?}", value, 1)
        }

        url = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(url, "")
        return url
}

func (r *Router) Routes() []*Route {
        return r.routes
}
