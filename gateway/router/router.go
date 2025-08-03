package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go-inventory-system/shared"
)

// Router handles routing requests to backend services
type Router struct {
	routes []shared.Route
	mux    *http.ServeMux
}

// NewRouter creates a new router with the given routes
func NewRouter(routes []shared.Route) *Router {
	router := &Router{
		routes: routes,
		mux:    http.NewServeMux(),
	}

	router.setupRoutes()
	return router
}

// setupRoutes configures the routing rules
func (r *Router) setupRoutes() {
	for _, route := range r.routes {
		backendURL, err := url.Parse(route.Backend)
		if err != nil {
			panic("Invalid backend URL: " + route.Backend)
		}

		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		proxy.ModifyResponse = r.modifyResponse

		// Create handler for this route
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Check if method is allowed
			if !r.isMethodAllowed(req.Method, route.Methods) {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Forward request to backend
			proxy.ServeHTTP(w, req)
		})

		// Register route
		r.mux.Handle(route.Path+"/", http.StripPrefix(route.Path, handler))
		r.mux.Handle(route.Path, handler)
	}

	// Health check endpoint
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// isMethodAllowed checks if the HTTP method is allowed for the route
func (r *Router) isMethodAllowed(method string, allowedMethods []string) bool {
	for _, allowed := range allowedMethods {
		if strings.EqualFold(method, allowed) {
			return true
		}
	}
	return false
}

// modifyResponse can be used to modify responses from backend services
func (r *Router) modifyResponse(resp *http.Response) error {
	// Add CORS headers
	resp.Header.Set("Access-Control-Allow-Origin", "*")
	resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	return nil
}
