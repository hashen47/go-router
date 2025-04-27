package grouter

import (
	"net/http"
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
