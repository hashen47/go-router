package main

import (
	"net/http"
	"testing"
)

func TestJoinPathAtFront(t *testing.T) {
	type TestCase struct {
		prefix      string
		pattern     string
		joinPattern string
	}

	testcases := []TestCase{
		{
			"/one/",
			"/two/three",
			"/one/two/three",
		},
		{
			"/one",
			"/two/three",
			"/one/two/three",
		},
		{
			"/one////////",
			"/two/three",
			"/one/two/three",
		},
		{
			"////one////",
			"/two/three",
			"/one/two/three",
		},
		{
			"////one////two/three//four",
			"/five/six",
			"/one/two/three/four/five/six",
		},
		{
			"one////two/three//four",
			"/five/six",
			"one/two/three/four/five/six",
		},
		{
			"one/two/three//four",
			"/five/six/s/{$}",
			"one/two/three/four/five/six/s/{$}",
		},
	}

	for _, tc := range testcases {
		route, _ := newRoute(tc.pattern, GET, func(w http.ResponseWriter, r *http.Request) {})
		route.joinPathAtFront(tc.prefix)

		if route.pattern != tc.joinPattern {
			t.Fatalf("ERROR\nPATTERN: %s\nPREFIX: %s\nEXPECT: %s\nREAL: %s\n", tc.pattern, tc.prefix, tc.joinPattern, route.pattern)
		}
		// t.Logf("Joined Pattern: %s\n", route.pattern)
	}
}

func TestNewRoute(t *testing.T) {
	type Testcase struct {
		pattern string
		rtype   RouteType
		err     error
	}

	testcases := []Testcase{
		{
			"/one",
			GET,
			nil,
		},
		{
			"/two",
			POST,
			nil,
		},
		{
			"/three",
			PUT,
			nil,
		},
		{
			"/four",
			PATCH,
			nil,
		},
		{
			"/four/five/{id}",
			DELETE,
			nil,
		},
		{
			"/four/five/{id}/",
			GET,
			nil,
		},
		{
			"/four/{name}/five/{id}/",
			POST,
			nil,
		},
		{
			"/four/{name}/five/{id}/{$}",
			PUT,
			nil,
		},
		{
			"/four/{name}/five/{id}/{$}",
			PATCH,
			nil,
		},
		{
			"/four/{$}",
			PATCH,
			nil,
		},
		{
			"/four/five/seven/eight/{$}/dfd",
			POST,
			nil,
		},
		{
			"/four/five/seven/eight/{$}/dfd",
			POST,
			nil,
		},
		{
			"four",
			DELETE,
			&RoutePatternInvalidError{
				&Route{pattern: "four", rtype: DELETE},
				"route should start with '/'",
			},
		},
		{
			"/four//five",
			GET,
			&RoutePatternInvalidError{
				&Route{pattern: "/four//five", rtype: GET},
				"same place cannot have more than single '/'",
			},
		},
		{
			"/four//five/six/seven",
			POST,
			&RoutePatternInvalidError{
				&Route{pattern: "/four//five/six/seven", rtype: POST},
				"same place cannot have more than single '/'",
			},
		},
		{
			"/four/five/six/seven/{}",
			PUT,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/six/seven/{}", rtype: PUT},
				"at offset 21: empty wildcard",
			},
		},
		{
			"/four/five/{}/seven/dfd",
			PATCH,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/{}/seven/dfd", rtype: PATCH},
				"at offset 11: empty wildcard",
			},
		},
		{
			"/four/five/six{id}/seven/dfd",
			DELETE,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/six{id}/seven/dfd", rtype: DELETE},
				"at offset 13: bad wildcard segment (must start with '{')",
			},
		},
		{
			"/four{id}/five/six/seven/dfd",
			GET,
			&RoutePatternInvalidError{
				&Route{pattern: "/four{id}/five/six/seven/dfd", rtype: GET},
				"at offset 4: bad wildcard segment (must start with '{')",
			},
		},
		{
			"/{id}four/five/six/seven/dfd",
			POST,
			&RoutePatternInvalidError{
				&Route{pattern: "/{id}four/five/six/seven/dfd", rtype: POST},
				"at offset 5: bad wildcard segment (must end with '}')",
			},
		},
		{
			"/four/five/{id}six/seven/dfd",
			PUT,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/{id}six/seven/dfd", rtype: PUT},
				"at offset 15: bad wildcard segment (must end with '}')",
			},
		},
		{
			"/four/five/{id}six/seven/dfd",
			PATCH,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/{id}six/seven/dfd", rtype: PATCH},
				"at offset 15: bad wildcard segment (must end with '}')",
			},
		},
		{
			"/four/five/{id}/seven/{id}/dfd",
			DELETE,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/five/{id}/seven/{id}/dfd", rtype: DELETE},
				"at offset 22: duplicate wildcard name 'id'",
			},
		},
		{
			"/four/{something}/five/seven/eight/{something}/dfd",
			GET,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/{something}/five/seven/eight/{something}/dfd", rtype: GET},
				"at offset 35: duplicate wildcard name 'something'",
			},
		},
		{
			"/{$}/four/{something}/five/seven/eight/{something}/dfd",
			POST,
			&RoutePatternInvalidError{
				&Route{pattern: "/{$}/four/{something}/five/seven/eight/{something}/dfd", rtype: POST},
				"at offset 1: {$} not at end",
			},
		},
		{
			"/four/{something}/five/{$}/seven/eight/{else}/dfd",
			PUT,
			&RoutePatternInvalidError{
				&Route{pattern: "/four/{something}/five/{$}/seven/eight/{else}/dfd", rtype: PUT},
				"at offset 23: {$} not at end",
			},
		},
	}

	for _, tc := range testcases {
		_, err := newRoute(tc.pattern, tc.rtype, func(w http.ResponseWriter, r *http.Request) {})

		if tc.err == nil && err != nil ||
			tc.err != nil && err == nil {
			t.Fatalf("ERROR\nPATTERN: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.pattern, tc.err, err)
		}

		if err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("ERROR\nPATTERN: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.pattern, tc.err, err)
			}
		}
	}
}
