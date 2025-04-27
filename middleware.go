package grouter

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(handler http.Handler) http.Handler

func createMiddlewareStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			next = m(next)
		}
		return next
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%-8v %-6s %s\n", time.Since(start), r.Method, r.URL.Path)
	})
}
