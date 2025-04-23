package main

import (
	"fmt"
)

func validateRoutePattern(route *Route, isRouteGroupPrefix bool) error {
	if len(route.pattern) == 0 {
		if isRouteGroupPrefix {
			return &RouteGroupPrefixInvalidError{route.pattern, "route pattern cannot be empty"}
		}
		return &RoutePatternInvalidError{route, "route pattern cannot be empty"}
	}

	if route.pattern[0] != '/' {
		if isRouteGroupPrefix {
			return &RouteGroupPrefixInvalidError{route.pattern, "route should start with '/'"}
		}
		return &RoutePatternInvalidError{route, "route should start with '/'"}
	}

	// validate wildcards
	curlyBraceIndex := -1
	wildcards := make(map[string]bool, 0)
	wildcard := make([]byte, 0)
	lastWildcardOffset := -1
	for i := 0; i < len(route.pattern); i++ {
		if i < len(route.pattern)-1 {
			if route.pattern[i] == '/' && route.pattern[i+1] == '/' {
				if isRouteGroupPrefix {
					return &RouteGroupPrefixInvalidError{route.pattern, "same place cannot have more than single '/'"}
				}
				return &RoutePatternInvalidError{route, "same place cannot have more than single '/'"}
			}
		}

		if i > 0 {
			if route.pattern[i] == '{' {
				if curlyBraceIndex != -1 {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must end with '}')", i),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must end with '}')", i),
					}
				}

				if route.pattern[i-1] != '/' {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must start with '{')", i-1),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must start with '{')", i-1),
					}
				}

				curlyBraceIndex = i
				continue
			}

			if route.pattern[i] == '}' {
				if curlyBraceIndex == -1 {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must start with '{')", i-1),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must start with '{')", i-1),
					}
				}

				if i != len(route.pattern)-1 {
					if route.pattern[i+1] != '/' {
						if isRouteGroupPrefix {
							return &RouteGroupPrefixInvalidError{
								prefix:   route.pattern,
								extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must end with '}')", i+1),
							}
						}
						return &RoutePatternInvalidError{
							route:    route,
							extraMsg: fmt.Sprintf("at offset %d: bad wildcard segment (must end with '}')", i+1),
						}
					}
				}

				if len(wildcard) == 0 {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: empty wildcard", curlyBraceIndex),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: empty wildcard", curlyBraceIndex),
					}
				}

				wildcardStr := string(wildcard)

				if _, ok := wildcards[wildcardStr]; ok {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: duplicate wildcard name '%s'", curlyBraceIndex, wildcardStr),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: duplicate wildcard name '%s'", curlyBraceIndex, wildcardStr),
					}
				}

				if lastWildcardOffset != -1 {
					if isRouteGroupPrefix {
						return &RouteGroupPrefixInvalidError{
							prefix:   route.pattern,
							extraMsg: fmt.Sprintf("at offset %d: {$} not at end", lastWildcardOffset),
						}
					}
					return &RoutePatternInvalidError{
						route:    route,
						extraMsg: fmt.Sprintf("at offset %d: {$} not at end", lastWildcardOffset),
					}
				}

				if wildcardStr == "$" {
					lastWildcardOffset = curlyBraceIndex
				}

				wildcards[wildcardStr] = true
				wildcard = make([]byte, 0)
				curlyBraceIndex = -1

				continue
			}

			if curlyBraceIndex != -1 {
				wildcard = append(wildcard, route.pattern[i])
			}
		}
	}

	return nil
}
