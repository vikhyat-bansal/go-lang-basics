package mux

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const varsContextKey contextKey = "vars"

// Router is a lightweight HTTP router compatible with basic gorilla/mux usage.
type Router struct {
	routes []*Route
}

// Route represents a route definition.
type Route struct {
	pattern string
	handler http.HandlerFunc
	methods map[string]struct{}
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{routes: make([]*Route, 0)}
}

// HandleFunc registers a path pattern with a handler function.
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Route {
	route := &Route{
		pattern: normalize(pattern),
		handler: handler,
		methods: map[string]struct{}{},
	}
	r.routes = append(r.routes, route)
	return route
}

// Methods limits a route to one or more HTTP methods.
func (rt *Route) Methods(methods ...string) *Route {
	for _, method := range methods {
		rt.methods[strings.ToUpper(strings.TrimSpace(method))] = struct{}{}
	}
	return rt
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestPath := normalize(req.URL.Path)
	pathMatched := false

	for _, route := range r.routes {
		vars, ok := match(route.pattern, requestPath)
		if !ok {
			continue
		}
		pathMatched = true

		if len(route.methods) > 0 {
			if _, exists := route.methods[req.Method]; !exists {
				continue
			}
		}

		ctx := context.WithValue(req.Context(), varsContextKey, vars)
		route.handler(w, req.WithContext(ctx))
		return
	}

	if pathMatched {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, req)
}

// Vars returns path variables for the current request.
func Vars(r *http.Request) map[string]string {
	vars, ok := r.Context().Value(varsContextKey).(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return vars
}

func normalize(path string) string {
	if path == "" {
		return "/"
	}
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func match(pattern, path string) (map[string]string, bool) {
	pattern = normalize(pattern)
	path = normalize(path)

	patternParts := split(pattern)
	pathParts := split(path)
	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	vars := make(map[string]string)
	for i := range patternParts {
		pp := patternParts[i]
		ap := pathParts[i]
		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(pp, "{"), "}")
			if key == "" {
				return nil, false
			}
			vars[key] = ap
			continue
		}
		if pp != ap {
			return nil, false
		}
	}

	return vars, true
}

func split(path string) []string {
	if path == "/" {
		return []string{}
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}
