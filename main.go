package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	. "web-app/router"
)

var router = Node(
	"/hello",  Leaf("GET", Handler(HelloWorld)),
	"/square", Leaf("GET", Handler(Square)),
	"/auth", WrapNode(
		Middlewares(Auth),
		"/square", Leaf(
			"GET",  Handler(Square),
			"POST", Handler(SquarePost),
		),
	),
)

// MIDDLEWARE
func Auth(f Handler) Handler {
	return func(r *http.Request) (int, any, error) {
		token := r.Header.Get("Authorization")
		if token == "" {
			return 401, nil, fmt.Errorf("use authorization header")
		}

		if token != "password" {
			return 401, nil, fmt.Errorf("bad password")
		}

		return f(r)
	}
}

// ENDPOINTS
func HelloWorld(r *http.Request) (int, any, error) {
	return 200, "Hello world", nil
}

// input type, parsed from GET args
// only strings allowed
type SquareQuery struct {
	Number string `json:"n"`
}

func Square(r *http.Request) (int, any, error) {
	var query SquareQuery
	err := ParseInput(r, &query)
	if err != nil {
		return 400, nil, err
	}

	n, err := strconv.Atoi(query.Number)
	return 200, n * n, err
}

// input type, parsed from POST JSON
type SquarePostQuery struct {
	Number int `json:"n"`
}

func SquarePost(r *http.Request) (int, any, error) {
	var query SquarePostQuery
	err := ParseInput(r, &query)
	if err != nil {
		return 400, nil, err
	}

	n := query.Number
	return 200, n * n, err
}

func main() {
	// load environment and workers
	port := flag.Int("p", 80, "Port to serve on")
	flag.Parse()

	// start web server
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", *port),
		BuildRouter(router),
	)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
