package grouter

import (
	"net/http"
	"strings"
)

type Router struct {
	addr              string
	routes            map[string]*Route
	collabRouteGroups map[int]bool
}

func initialize(addr string) (*Router, error) {
	// TODO: apply some validations to check router address is correct
	addr = strings.Trim(addr, " ")
	if len(addr) == 0 {
		return nil, &RouterAddrEmptyError{}
	}

	router := &Router{
		addr:              addr,
		routes:            make(map[string]*Route, 0),
		collabRouteGroups: make(map[int]bool, 0),
	}

	return router, nil
}

func Init(addr string) *Router {
	router, err := initialize(addr)
	if err != nil {
		printAndExit(err, 1)
	}
	return router
}

func (r *Router) putGroup(rg *RouteGroup) error {
	if _, ok := r.collabRouteGroups[rg.id]; ok {
		return &CircularRouteGroupAddingError{rg.id, rg.prefix}
	}

	if len(rg.routes) == 0 {
		return &RouterEmptyError{}
	}

	if err := rg.applyMiddlewares(); err != nil {
		return err
	}

	for _, route := range rg.routes {
		routeTypeStr, _ := getRouteTypeStr(route)
		route.pattern = routeTypeStr + " " + route.pattern
		if _, ok := r.routes[route.pattern]; ok {
			return &DuplicateRouteError{route}
		}
		r.routes[route.pattern] = route
	}

	for key, val := range rg.collabRouteGroups {
		r.collabRouteGroups[key] = val
	}

	return nil
}

func (r *Router) AddGroup(rg *RouteGroup) {
	if err := r.putGroup(rg); err != nil {
		printAndExit(err, 1)
	}
}

func (r *Router) getServeMux() (http.Handler, error) {
	if len(r.routes) == 0 {
		return nil, &RouterEmptyError{}
	}

	server := http.NewServeMux()

	for _, route := range r.routes {
		middlewareStack := createMiddlewareStack(route.middlewares...)
		server.Handle(route.pattern, middlewareStack(http.HandlerFunc(route.handler)))
	}

	return server, nil
}

func (r *Router) ListenAndServe() {
	handler, err := r.getServeMux()
	if err != nil {
		printAndExit(err, 1)
	}
	printAndExit(http.ListenAndServe(r.addr, handler), 1)
}
