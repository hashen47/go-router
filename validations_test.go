package main

import (
	"testing"
)

func TestValidateRoutePattern(t *testing.T) {
	type TestCase struct {
		route *Route
		err   error
	}

	testcases := []TestCase{
		{
			route: &Route{pattern: "/one/two/three/four/five"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/about"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/contact"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user/123"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/post/456/comments"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/product/shirt"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/search?q=golang"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/filter?category=books&price=10-20"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/login?redirect=/dashboard"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/api/users/42/orders"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/api/users/42/orders/99"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user-profile/42"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user_profile/42"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user_profile?q=22&q=22"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user/😀"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/four/{$}/five/seven/eight/dfd"},
			err:   nil,
		},
		{
			route: &Route{pattern: "/user/{}"},
			err: &RoutePatternInvalidError{
				route:    &Route{pattern: "/user/{}"},
				extraMsg: "at offset 6: empty wildcard",
			},
		},
		{
			route: &Route{pattern: "/user/{$}/{id}"},
			err: &RoutePatternInvalidError{
				route:    &Route{pattern: "/user/{$}/{id}"},
				extraMsg: "at offset 6: {$} not at end",
			},
		},
	}

	for _, tc := range testcases {
		err := validateRoutePattern(tc.route, false)
		if err == nil && tc.err != nil ||
			err != nil && tc.err == nil {
			t.Fatalf("\nERROR:\npattern: \"%s\",\nerror(expect): \"%v\",\nerror(real)  : \"%v\"\n", tc.route.pattern, tc.err, err)
		}

		if err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("\nERROR:\npattern: \"%s\",\nerror(expect): \"%v\",\nerror(real)  : \"%v\"\n", tc.route.pattern, tc.err, err)
			}
		}
	}
}
