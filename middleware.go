package grouter

import (
	"net/http"
)

type Middleware func(handler http.Handler) http.Handler

func createMiddlewareStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(middlewares); i++ {
			m := middlewares[i]
			next = m(next)
		}
		return next
	}
}
