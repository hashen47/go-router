package main

import (
	"net/http"
	"strings"
)

type RouteGroup struct {
	id                int
	prefix            string
	routes            map[string]*Route
	collabRouteGroups map[int]bool
	middlewares       []Middleware
}

func NewRouteGroup(prefix string) *RouteGroup {
	prefix = strings.Trim(prefix, " ")

	if len(prefix) != 0 {
		route, _ := newRoute(prefix, GET, func(w http.ResponseWriter, r *http.Request) {})
		if err := validateRoutePattern(route, true); err != nil {
			printAndExit(err, 1)
		}
	}

	route := &RouteGroup{
		id:                generateUniqueNumber(),
		prefix:            prefix,
		routes:            make(map[string]*Route, 0),
		middlewares:       make([]Middleware, 0),
		collabRouteGroups: make(map[int]bool, 0),
	}
	route.collabRouteGroups[route.id] = true

	return route
}

func (g *RouteGroup) putRoute(pattern string, rtype RouteType, handler Handler) (*Route, error) {
	route, err := newRoute(pattern, rtype, handler)
	if err != nil {
		return route, err
	}

	if g.prefix != "" {
		route.joinPathAtFront(g.prefix)
	}

	if r, ok := g.routes[route.pattern]; ok {
		if r.rtype == route.rtype {
			return route, &DuplicateRouteError{route}
		}
	}

	if err := validateRoutePattern(route, true); err != nil {
		return route, err
	}

	return route, nil
}

func (g *RouteGroup) AddRoute(pattern string, rtype RouteType, handler Handler) *Route {
	route, err := g.putRoute(pattern, rtype, handler)
	if err != nil {
		printAndExit(err, 1)
	}
	g.routes[route.pattern] = route
	return route
}

func (g *RouteGroup) putRouteGroup(rg *RouteGroup) error {
	if _, ok := g.collabRouteGroups[rg.id]; ok {
		return &CircularRouteGroupAddingError{g.id, g.prefix}
	}

	if _, ok := rg.collabRouteGroups[g.id]; ok {
		return &CircularRouteGroupAddingError{g.id, g.prefix}
	}

	if err := rg.applyMiddlewares(); err != nil {
		return err
	}

	for routePattern, route := range rg.routes {
		if r, ok := g.routes[routePattern]; ok {
			if g.prefix != "" && r.rtype == route.rtype {
				return &DuplicateRouteError{route}
			}
		}
		if _, err := g.putRoute(route.pattern, route.rtype, route.handler); err != nil {
			return err
		}
		r := g.AddRoute(route.pattern, route.rtype, route.handler)
		r.AddMiddlewares(route.middlewares...)
	}

	rg.routes = make(map[string]*Route, 0)

	tempMap := make(map[int]bool, 0)

	for key, val := range g.collabRouteGroups {
		tempMap[key] = val
	}

	for key, val := range rg.collabRouteGroups {
		tempMap[key] = val
	}

	g.collabRouteGroups = make(map[int]bool, 0)
	rg.collabRouteGroups = make(map[int]bool, 0)

	for key, val := range tempMap {
		g.collabRouteGroups[key] = val
		rg.collabRouteGroups[key] = val
	}

	return nil
}

func (g *RouteGroup) AddRouteGroup(rg *RouteGroup) {
	if err := g.putRouteGroup(rg); err != nil {
		printAndExit(err, 1)
	}
}

func (g *RouteGroup) AddMiddlewaresToGroup(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouteGroup) applyMiddlewares() error {
	if len(g.routes) == 0 {
		return &RouteGroupEmptyError{g.id, g.prefix}
	}

	for _, route := range g.routes {
		route.AddMiddlewares(g.middlewares...)
	}

	return nil
}
