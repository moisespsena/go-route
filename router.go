//
// Package chi is a small, idiomatic and composable router for building HTTP services.
//
// chi requires Go 1.7 or newer.
//
// Example:
//  package main
//
//  import (
//  	"net/http"
//
//  	"github.com/go-chi/chi"
//  	"github.com/go-chi/chi/middleware"
//  )
//
//  func main() {
//  	r := chi.NewRouter()
//  	r.Use(middleware.Logger)
//  	r.Use(middleware.Recoverer)
//
//  	r.Get("/", func(w http.NewResponseWriter, r *http.request) {
//  		w.Write([]byte("root."))
//  	})
//
//  	http.ListenAndServe(":3333", r)
//  }
//
// See github.com/go-chi/chi/_examples/ for more in-depth examples.
//
package xroute

import "net/http"

// NewRouter returns a new Mux object that implements the Router interface.
func NewRouter() *Mux {
	return NewMux()
}

type HandlerFinder interface {
	FindHandler(method, path string, header ...http.Header) ContextHandler
}

// Router consisting of the core routing methods used by chi's Mux,
// using only the standard net/http.
type Router interface {
	Handler
	Routes
	HandlerFinder

	Prefix() string
	SetPrefix(p string)

	SetRouteHandler(handler ContextHandlerFunc)
	GetRouteHandler() ContextHandlerFunc

	IsArgSet() bool
	SetArg(arg interface{})
	Arg() interface{}
	ClearArg()

	Intersept(interseptors ...interface{})
	GetInterseptor(name string) *Middleware

	HandlerIntersept(interseptors ...interface{})
	GetHandlerInterseptor(name string) *Middleware

	// Use appends one of more middlewares onto the Router stack.
	Use(middlewares ...interface{})

	// Return middleware by name or `nil`
	GetMiddleware(name string) *Middleware

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...interface{}) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a Sub-Router along a `pattern` string.
	Route(pattern string, fn func(r Router)) Router

	// Mount attaches another interface{} along ./pattern/*
	Mount(pattern string, h interface{})

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h interface{})

	// Method and add routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h interface{})

	// MethodT adds the route `pattern` that matches `method` http method to
	// execute the `handler` Handler.
	MethodT(method MethodType, pattern string, handler interface{})

	// HTTP-method routing along `pattern`
	HandleMethod(method string, pattern string, handler interface{})
	HandleM(method MethodType, pattern string, handler interface{})
	Connect(pattern string, h interface{})
	Delete(pattern string, h interface{})
	Get(pattern string, h interface{})
	Head(pattern string, h interface{})
	Options(pattern string, h interface{})
	Patch(pattern string, h interface{})
	Post(pattern string, h interface{})
	Put(pattern string, h interface{})
	Trace(pattern string, h interface{})

	Headers(headers http.Header, f func(r Router))
	Api(f func(r Router))

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h interface{})

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h interface{})

	Overrides(f func(r Router))
}

// Routes interface adds two methods for router traversal, which is also
// used by the `docgen` subpackage to generation documentation for Routers.
type Routes interface {
	// Routes returns the routing tree in an easily traversable structure.
	Routes() []Route

	// Middlewares returns the list of middlewares in use by the router.
	Middlewares() Middlewares

	// Match searches the routing tree for a handler that matches
	// the method/path - similar to routing a http request, but without
	// executing the handler thereafter.
	Match(rctx *RouteContext, method, path string) bool
}
