package router

import (
	"net/http"
)

// Middleware has to match this function definition
type Middleware func(w http.ResponseWriter, r *http.Request, c *Context) bool

// A Group is needed to call post, get, etc.
// Routes registered with a group will use the middleware registered with that group
type Group struct {
	middleware []Middleware
}

// NewGroup returns a *Group object
func NewGroup() *Group {
	group := &Group{
		middleware: []Middleware{},
	}
	return group
}

// Use registers a middleware with the group that calls it
func (g *Group) Use(m Middleware) {
	g.middleware = append(g.middleware, m)
}

// runMiddleware is called by the route object in the match method right before the callback is called
func (g *Group) runMiddleware(w http.ResponseWriter, req *http.Request, c *Context) bool {
	for _, m := range g.middleware {
		ok := m(w, req, c)
		if !ok {
			return false
		}
	}
	return true
}
