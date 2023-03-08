package router

import "net/http"

// BuildDefaultHTTPRouter creates a default HTTP router.
func BuildDefaultHTTPRouter() http.Handler {
	return http.NotFoundHandler()
}
