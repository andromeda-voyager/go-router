package router

import (
	"net/http"
)

// root route for all nested routes
var routes *route

func init() {
	routes = &route{name: "/", nestedRoutes: make(map[string]*route)}
}

// routeCallbackFunc defines the method definition for a route callback function
type routeCallbackFunc func(w http.ResponseWriter, r *http.Request, c *Context)

func setHeaders(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// (*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

// Post registers the provided callback and group that calls it with a POST route
func (g *Group) Post(path string, callback routeCallbackFunc) {
	routes.build(splitPath(path, "POST"), callback, g)
}

// Delete registers the provided callback and group that calls it with a DELETE route
func (g *Group) Delete(path string, callback routeCallbackFunc) {
	routes.build(splitPath(path, "DELETE"), callback, g)
}

// Get registers the provided callback and group that calls it with a GET route
func (g *Group) Get(path string, callback routeCallbackFunc) {
	routes.build(splitPath(path, "GET"), callback, g)
}

// Put registers the provided callback and group that calls it with a PUT route
func (g *Group) Put(path string, callback routeCallbackFunc) {
	routes.build(splitPath(path, "PUT"), callback, g)
}

// Handler is called by http.HandleFunc. The match function is called on the root route object
// to find the route assoicated with the req url path.
func Handler(w http.ResponseWriter, req *http.Request) {
	setHeaders(&w, req)
	if req.Method != "OPTIONS" {
		c := &Context{Keys: make(map[string]interface{})}
		routes.match(splitPath(req.URL.Path, req.Method), w, req, c)
	}
}

func PrintRoutes() {
	routes.print()
}
