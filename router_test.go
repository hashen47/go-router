package grouter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPutGroup(t *testing.T) {
	// TODO: implement more testcases to test this

	router := Init(":8080")

	g1 := NewRouteGroup("/v1")
	g2 := NewRouteGroup("/book")
	g3 := NewRouteGroup("/page")
	g4 := NewRouteGroup("/chapter")

	g2.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})
	g3.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})
	g4.AddRoute("/", GET, func(w http.ResponseWriter, r *http.Request) {})

	g3.AddRouteGroup(g4)
	g2.AddRouteGroup(g3)
	g1.AddRouteGroup(g2)

	err := router.putGroup(g1)
	if err != nil {
		t.Fatalf("ERROR")
	}
	// t.Log(router.collabRouteGroups)

	err = router.putGroup(g2)
	if err == nil {
		t.Fatalf("ERROR")
	}
}

func m1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prevVal string

		if v, ok := r.Context().Value("m").(string); ok {
			prevVal = v
		}

		newCtx := context.WithValue(r.Context(), "m", prevVal+"m1")
		newRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, newRequest)
	})
}

func m2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prevVal string

		if v, ok := r.Context().Value("m").(string); ok {
			prevVal = v
		}

		newCtx := context.WithValue(r.Context(), "m", prevVal+"m2")
		newRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, newRequest)
	})
}

func m3(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prevVal string

		if v, ok := r.Context().Value("m").(string); ok {
			prevVal = v
		}

		newCtx := context.WithValue(r.Context(), "m", prevVal+"m3")
		newRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, newRequest)
	})
}

func m4(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prevVal string

		if v, ok := r.Context().Value("m").(string); ok {
			prevVal = v
		}

		newCtx := context.WithValue(r.Context(), "m", prevVal+"m4")
		newRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, newRequest)
	})
}

func m5(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prevVal string

		if v, ok := r.Context().Value("m").(string); ok {
			prevVal = v
		}

		newCtx := context.WithValue(r.Context(), "m", prevVal+"m5")
		newRequest := r.WithContext(newCtx)
		next.ServeHTTP(w, newRequest)
	})
}

type RouteData struct {
	pattern     string
	rtype       RouteType
	handler     Handler
	middlewares []Middleware
}

type RouteGroupData struct {
	prefix      string
	middlewares []Middleware
	routes      []RouteData
	groups      []RouteGroupData
}

type Response struct {
	pattern     string
	data        string
	rtype       RouteType
	isHeaderSet bool
	headers     map[string]string
	statusCode  int
	err         error
}

type Testcase struct {
	groups    []RouteGroupData
	responses []Response
}

func TestGetServeMux(t *testing.T) {
	// TODO: add some serious test cases too.. (fail testcases)

	testcases := []Testcase{
		{
			groups: []RouteGroupData{
				{
					prefix:      "/v1/",
					middlewares: []Middleware{m1},
					routes: []RouteData{
						{
							pattern: "/one",
							rtype:   GET,
							handler: func(w http.ResponseWriter, r *http.Request) {
								val := ""

								if v, ok := r.Context().Value("m").(string); ok {
									val = v
								}

								w.Header().Add("x-header", "x-header-value")

								w.WriteHeader(http.StatusInternalServerError)

								fmt.Fprintf(w, val+"/v1/one")
							},
							middlewares: []Middleware{m2},
						},
					},
					groups: []RouteGroupData{
						{
							prefix:      "/two",
							middlewares: []Middleware{m2, m3},
							routes: []RouteData{
								{
									pattern: "/three",
									rtype:   GET,
									handler: func(w http.ResponseWriter, r *http.Request) {
										val := ""

										if v, ok := r.Context().Value("m").(string); ok {
											val = v
										}

										w.Header().Add("header-one", "val1")
										w.Header().Add("header-two", "val2")
										w.Header().Add("Content-Type", "application/json")

										w.WriteHeader(http.StatusPartialContent)

										fmt.Fprintf(w, val+"/v1/two/three")
									},
									middlewares: []Middleware{m4},
								},
							},
							groups: []RouteGroupData{
								{
									prefix: "/four",
									groups: []RouteGroupData{},
									routes: []RouteData{
										{
											pattern:     "/five",
											rtype:       GET,
											middlewares: []Middleware{},
											handler: func(w http.ResponseWriter, r *http.Request) {
												val := ""

												if v, ok := r.Context().Value("m").(string); ok {
													val = v
												}

												w.WriteHeader(http.StatusBadRequest)

												fmt.Fprintf(w, val+"/v1/two/four/five")
											},
										},
									},
								},
							},
						},
					},
				},
			},
			responses: []Response{
				{
					pattern:     "/v1/one",
					data:        "m2m1/v1/one",
					rtype:       GET,
					isHeaderSet: false,
					statusCode:  http.StatusInternalServerError,
					headers: map[string]string{
						"x-header": "x-header-value",
					},
					err: nil,
				},
				{
					pattern:     "/v1/two/three",
					data:        "m4m3m2m1/v1/two/three",
					rtype:       GET,
					isHeaderSet: false,
					statusCode:  http.StatusPartialContent,
					headers: map[string]string{
						"header-one":   "val1",
						"header-two":   "val2",
						"Content-Type": "application/json",
					},
					err: nil,
				},
				{
					pattern:     "/v1/two/four/five",
					data:        "m3m2m1/v1/two/four/five",
					rtype:       GET,
					isHeaderSet: false,
					statusCode:  http.StatusBadRequest,
					err:         nil,
				},
			},
		},
		{
			groups: []RouteGroupData{
				{
					prefix:      "",
					middlewares: []Middleware{m1},
					routes: []RouteData{
						{
							pattern:     "/one",
							rtype:       POST,
							middlewares: []Middleware{m2, m3},
							handler: func(w http.ResponseWriter, r *http.Request) {
								val := ""

								if v, ok := r.Context().Value("m").(string); ok {
									val = v
								}

								w.Header().Add("test-header-1", "value-test-one")
								w.Header().Add("test-header-2", "value-test-two")
								w.Header().Add("another-header", "test-value")

								w.WriteHeader(http.StatusCreated)

								fmt.Fprintf(w, val+"/one")
							},
						},
					},
					groups: []RouteGroupData{
						{
							prefix:      "/two",
							middlewares: []Middleware{m4},
							routes: []RouteData{
								{
									pattern:     "/three",
									rtype:       DELETE,
									middlewares: []Middleware{m5},
									handler: func(w http.ResponseWriter, r *http.Request) {
										val := ""

										if v, ok := r.Context().Value("m").(string); ok {
											val = v
										}

										w.Header().Add("test-header-1", "value-test-one")
										w.Header().Add("test-header-2", "value-test-two")
										w.Header().Add("another-header", "test-value")
										w.Header().Add("another-header-two", "test-value")
										w.Header().Add("a", "test-value")

										w.WriteHeader(http.StatusOK)

										fmt.Fprintf(w, val+"/two/three")
									},
								},
							},
						},
					},
				},
			},
			responses: []Response{
				{
					pattern:     "/one",
					data:        "m3m2m1/one",
					rtype:       POST,
					isHeaderSet: false,
					statusCode:  http.StatusCreated,
					headers: map[string]string{
						"test-header-1":  "value-test-one",
						"test-header-2":  "value-test-two",
						"another-header": "test-value",
					},
					err: nil,
				},
				{
					pattern:     "/two/three",
					data:        "m5m4m1/two/three",
					rtype:       DELETE,
					isHeaderSet: false,
					statusCode:  http.StatusOK,
					headers: map[string]string{
						"test-header-1":      "value-test-one",
						"test-header-2":      "value-test-two",
						"another-header":     "test-value",
						"another-header-two": "test-value",
						"a":                  "test-value",
					},
					err: nil,
				},
			},
		},
	}

	for _, tc := range testcases {
		router := Init(":8080")

		for _, groupData := range tc.groups {
			group := depth(groupData)
			router.AddGroup(group)
		}

		handler, _ := router.getServeMux()

		server := httptest.NewServer(handler)

		for _, response := range tc.responses {
			var resp *http.Response
			var err error
			var request *http.Request

			method := http.MethodGet

			switch response.rtype {
			case GET:
				method = http.MethodGet
			case POST:
				method = http.MethodPost
			case PUT:
				method = http.MethodPut
			case PATCH:
				method = http.MethodPatch
			case DELETE:
				method = http.MethodDelete
			}

			request, err = http.NewRequest(method, server.URL+response.pattern, nil)
			client := http.DefaultClient

			resp, err = client.Do(request)

			defer resp.Body.Close()

			if err != nil && response.err == nil ||
				err == nil && response.err != nil {
				t.Fatalf("\nERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.err, err)
			}

			if err != nil {
				if err.Error() != response.err.Error() {
					t.Fatalf("\nERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.err, err)
				}
			}

			if resp.StatusCode != response.statusCode {
				t.Fatalf("\nSTATUS CODE NOT MATCH ERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.statusCode, resp.StatusCode)
			}

			for header, val := range response.headers {
				v := resp.Header.Get(header)
				if val != v {
					t.Fatalf("\nHEADER VALUE NOT MATCH ERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, val, v)
				}
			}

			data, err := io.ReadAll(resp.Body)

			if err != nil && response.err == nil ||
				err == nil && response.err != nil {
				t.Fatalf("\nERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.err, err)
			}

			if err != nil {
				if err.Error() != response.err.Error() {
					t.Fatalf("\nERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.err, err)
				}
			}

			if string(data) != response.data {
				t.Fatalf("\nERROR\nPATTERN: %s\nEXPECT: %v\nREAL: %v\n", server.URL+response.pattern, response.data, string(data))
			}
		}
	}
}

func depth(groupData RouteGroupData) *RouteGroup {
	group := NewRouteGroup(groupData.prefix)
	group.AddMiddlewaresToGroup(groupData.middlewares...)

	for _, routeData := range groupData.routes {
		route := group.AddRoute(routeData.pattern, routeData.rtype, routeData.handler)
		route.AddMiddlewares(routeData.middlewares...)
	}

	for _, subGroupData := range groupData.groups {
		subGroup := depth(subGroupData)
		group.AddRouteGroup(subGroup)
	}

	return group
}
