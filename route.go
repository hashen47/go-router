package grouter

import (
	"fmt"
	"strings"
)

type RouteType int

const (
	GET RouteType = iota
	POST
	PUT
	PATCH
	DELETE
)

type Route struct {
	id          int
	pattern     string
	rtype       RouteType
	handler     Handler
	middlewares []Middleware
}

func newRoute(pattern string, rtype RouteType, handler Handler) (*Route, error) {
	route := &Route{
		id:          generateUniqueNumber(),
		pattern:     strings.Trim(pattern, " "),
		rtype:       rtype,
		handler:     handler,
		middlewares: make([]Middleware, 0),
	}

	if _, err := getRouteTypeStr(route); err != nil {
		return route, err
	}

	if err := validateRoutePattern(route, false); err != nil {
		return route, err
	}

	return route, nil
}

func (r *Route) AddMiddlewares(middlewares ...Middleware) {
	newMiddlewares := make([]Middleware, 0)
	newMiddlewares = append(newMiddlewares, middlewares...)
	newMiddlewares = append(newMiddlewares, r.middlewares...)
	r.middlewares = newMiddlewares
}

func (r *Route) joinPathAtFront(path string) {
	pattern := path + r.pattern
	updatedPattern := make([]byte, 0)

	isForwardSlashFound := false
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '/' {
			if isForwardSlashFound {
				continue
			}
			isForwardSlashFound = true
		} else {
			isForwardSlashFound = false
		}
		updatedPattern = append(updatedPattern, pattern[i])
	}

	r.pattern = string(updatedPattern)
}

func (r *Route) String() string {
	routeTypeStr, _ := getRouteTypeStr(r)
	return fmt.Sprintf("\n------------------------\nPATTERN: %s\nTYPE: %s\n------------------------\n", r.pattern, routeTypeStr)
}
