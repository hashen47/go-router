package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPutRoute(t *testing.T) {
	// main target is check route duplication error, because route base other errors check in other testcases
	type RouteData struct {
		pattern string
		rtype   RouteType
		handler Handler
	}

	type Testcase struct {
		prefix string
		routes []RouteData
		err    error
	}

	testcases := []*Testcase{
		{
			"",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			nil,
		},
		{
			"/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			nil,
		},
		{
			"/v1/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			nil,
		},
		{
			"/v1/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", POST, func(w http.ResponseWriter, r *http.Request) {}},
			},
			nil,
		},
		{
			"/v2/book/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			nil,
		},
		{
			"/one",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
			},
			&DuplicateRouteError{
				&Route{pattern: "/one/one", rtype: GET, handler: func(w http.ResponseWriter, r *http.Request) {}},
			},
		},
		{
			"/v1/",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
			},
			&DuplicateRouteError{
				&Route{pattern: "/v1/one", rtype: GET, handler: func(w http.ResponseWriter, r *http.Request) {}},
			},
		},
	}

	for _, tc := range testcases {
		g := NewRouteGroup(tc.prefix)
		var realErr error = nil

		for _, routeData := range tc.routes {
			_, err := g.putRoute(routeData.pattern, routeData.rtype, routeData.handler)
			// t.Log(route)
			if err != nil {
				realErr = err
				continue
			}
			// t.Log(route)
			g.AddRoute(routeData.pattern, routeData.rtype, routeData.handler)
		}

		if tc.err == nil && realErr != nil ||
			tc.err != nil && realErr == nil {
			t.Fatalf("ERROR\nPREFIX: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.prefix, tc.err, realErr)
		}

		if realErr != nil {
			if tc.err.Error() != realErr.Error() {
				t.Fatalf("ERROR\nPREFIX: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.prefix, tc.err, realErr)
			}
		}
	}
}

func TestPutRouteGroup(t *testing.T) {
	// TODO: create more test cases for this method
	g1 := NewRouteGroup("/v1")

	g2 := NewRouteGroup("/book")
	g3 := NewRouteGroup("/page")
	g4 := NewRouteGroup("/chapter")

	g2.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})
	g3.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})
	g4.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})

	var err error

	err = g3.putRouteGroup(g4)
	if err != nil {
		t.Fatalf("ERROR")
	}

	err = g2.putRouteGroup(g3)
	if err != nil {
		t.Fatalf("ERROR")
	}

	err = g1.putRouteGroup(g2)
	if err != nil {
		t.Fatalf("ERROR")
	}

	err = g1.putRouteGroup(g3)
	if err == nil {
		t.Fatalf("ERROR")
	}

	// for _, route := range g1.routes {
	// 	t.Log(route)
	// }
}

func MiddlewareOne(next http.Handler) http.Handler {
	fmt.Println("MIDDLEWARE ONE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func MiddlewareTwo(next http.Handler) http.Handler {
	fmt.Println("MIDDLEWARE TWO")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func MiddlewareThree(next http.Handler) http.Handler {
	fmt.Println("MIDDLEWARE THREE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func MiddlewareFour(next http.Handler) http.Handler {
	fmt.Println("MIDDLEWARE FOUR")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func MiddlewareFive(next http.Handler) http.Handler {
	fmt.Println("MIDDLEWARE FIVE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestApplyMiddlewares(t *testing.T) {
	type RouteData struct {
		pattern string
		rtype   RouteType
		handler Handler
	}

	type Testcase struct {
		id          int
		prefix      string
		routes      []RouteData
		middlewares []Middleware
		err         error
	}

	testcases := []*Testcase{
		{
			1,
			"",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
			},
			nil,
		},
		{
			2,
			"/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
				MiddlewareThree,
			},
			nil,
		},
		{
			3,
			"/v1/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
				MiddlewareThree,
				MiddlewareFour,
			},
			nil,
		},
		{
			3,
			"/v1/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{},
			nil,
		},
		{
			4,
			"/v1/page",
			[]RouteData{
				{"/one", GET, func(w http.ResponseWriter, r *http.Request) {}},
				{"/two", POST, func(w http.ResponseWriter, r *http.Request) {}},
				{"/three", PUT, func(w http.ResponseWriter, r *http.Request) {}},
				{"/four", PATCH, func(w http.ResponseWriter, r *http.Request) {}},
				{"/five", DELETE, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
				MiddlewareThree,
				MiddlewareFour,
				MiddlewareFive,
			},
			nil,
		},
		{
			5,
			"/v2/book/page",
			[]RouteData{
				{"/one/two/three/four/five", GET, func(w http.ResponseWriter, r *http.Request) {}},
			},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
				MiddlewareFour,
				MiddlewareFive,
			},
			nil,
		},
		{
			6,
			"/one",
			[]RouteData{},
			[]Middleware{
				MiddlewareOne,
				MiddlewareFour,
				MiddlewareFive,
			},
			&RouteGroupEmptyError{6, "/one"},
		},
		{
			7,
			"/v1/",
			[]RouteData{},
			[]Middleware{
				MiddlewareOne,
				MiddlewareTwo,
				MiddlewareThree,
				MiddlewareFour,
				MiddlewareFive,
			},
			&RouteGroupEmptyError{7, "/v1/"},
		},
	}

	for _, tc := range testcases {
		g := NewRouteGroup(tc.prefix)
		g.id = tc.id

		for _, routeData := range tc.routes {
			g.AddRoute(routeData.pattern, routeData.rtype, routeData.handler)
		}

		g.AddMiddlewaresToGroup(tc.middlewares...)
		err := g.applyMiddlewares()

		if tc.err == nil && err != nil ||
			tc.err != nil && err == nil {
			t.Fatalf("ERROR\nPREFIX: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.prefix, tc.err, err)
		}

		if err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("ERROR\nPREFIX: %s\nEXPECT:\n%v\nREAL:\n%v\n", tc.prefix, tc.err, err)
			}
		}
	}
}
