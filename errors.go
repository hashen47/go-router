package grouter

import (
	"fmt"
)

// Route Errors
type InvalidRouteTypeError struct {
	route *Route
}

func (e *InvalidRouteTypeError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: Invalid Route Type\nTYPE: %d\nPATTERN: %s\n----------------------------", e.route.rtype, e.route.pattern)
}

type RoutePatternInvalidError struct {
	route    *Route
	extraMsg string
}

func (e *RoutePatternInvalidError) Error() string {
	routeTypeStr, _ := getRouteTypeStr(e.route)
	return fmt.Sprintf("----------------------------\nERROR: Invalid Route Pattern\nTYPE: %s\nPATTERN: %s\nMESSAGE: %s\n----------------------------", routeTypeStr, e.route.pattern, e.extraMsg)
}

type DuplicateRouteError struct {
	route *Route
}

func (e *DuplicateRouteError) Error() string {
	routeTypeStr, _ := getRouteTypeStr(e.route)
	return fmt.Sprintf("----------------------------\nERROR: Route is duplicated\nTYPE: %s\nPATTERN: %s\n----------------------------", routeTypeStr, e.route.pattern)
}

// Route Group Errors
type RouteGroupPrefixInvalidError struct {
	prefix   string
	extraMsg string
}

func (e *RouteGroupPrefixInvalidError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: Prefix is invalid\nPREFIX: %s\nMESSAGE: %s\n----------------------------", e.prefix, e.extraMsg)
}

type CircularRouteGroupAddingError struct {
	id     int
	prefix string
}

func (e *CircularRouteGroupAddingError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: RouteGroups are circulary added\nID: %d\nPREFIX: %s\n----------------------------", e.id, e.prefix)
}

type RouteGroupEmptyError struct {
	id     int
	prefix string
}

func (e *RouteGroupEmptyError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: RouteGroup hasn't contain any route\nID: %d\nPREFIX: %s\n----------------------------", e.id, e.prefix)
}

// Router Errors
type RouterAddrEmptyError struct{}

func (e *RouterAddrEmptyError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: Router address cannot be empty\n----------------------------")
}

type RouterEmptyError struct{}

func (e *RouterEmptyError) Error() string {
	return fmt.Sprintf("----------------------------\nERROR: Router hasn't contain any route\n----------------------------")
}
