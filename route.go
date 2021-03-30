package router

import (
	"fmt"
	"net/http"
	"strings"
)

// route is used to build nested routes based on the url path registered with
// a route. It has fields for the callback function, group, and path parameters
// used with a route.
type route struct {
	name         string
	paramName    string
	isIntParam   bool
	isPathParam  bool
	nestedRoutes map[string]*route
	callback     routeCallbackFunc
	group        *Group
}

// splitPath takes a url path and returns the segments of that path
// split by '/'.
func splitPath(path, method string) []string {
	nestedPaths := strings.Split(path, "/")
	nestedPaths[0] = method
	return nestedPaths
}

// build creates the nested routes needed to map a url path with a
// specific callback and the group associated with that route.
func (r *route) build(pathSegments []string, fn routeCallbackFunc, g *Group) {
	if len(pathSegments) == 0 {
		r.callback = fn
		r.group = g
	} else {
		nestedRoute, ok := r.nestedRoutes[pathSegments[0]]
		if strings.HasPrefix(pathSegments[0], ":") {
			nestedRoute, ok = r.nestedRoutes["param"]
		}
		if ok {
			nestedRoute.build(pathSegments[1:], fn, g)
		} else {
			nestedRoute := &route{
				name:         pathSegments[0],
				nestedRoutes: make(map[string]*route),
			}
			if strings.HasPrefix(nestedRoute.name, ":") {
				nestedRoute.makePathParam()
			}

			nestedRoute.build(pathSegments[1:], fn, g)
			r.nestedRoutes[nestedRoute.name] = nestedRoute
		}
	}
}

func (r *route) print() {
	for x, nr := range r.nestedRoutes {
		if nr.callback != nil {
			if x == "param" {
				fmt.Println(nr.paramName)
			} else {
				fmt.Println(x)
			}
			for s := range nr.nestedRoutes {
				fmt.Println("---", s)
			}
		}
		nr.print()
	}
}

// match maps a url path to a registered route. If a route is found for
// a path, the middleware used with the group that created that route is
// run. If the middleware returns an ok to continue, then the callback
// registered with that route is run.
func (r *route) match(pathSegments []string, w http.ResponseWriter, req *http.Request, c *Context) {
	if len(pathSegments) == 0 {
		ok := r.group.runMiddleware(w, req, c)
		if ok {
			if r.callback != nil {
				r.callback(w, req, c)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	} else {
		nestedRoute, ok := r.getNestedRoute(pathSegments[0])
		if ok {
			if nestedRoute.isPathParam {
				ok := c.addPathParam(nestedRoute, pathSegments[0])
				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			nestedRoute.match(pathSegments[1:], w, req, c)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// getNestedRoute returns the nested route mapped to the provided path segment.
// If a the pathSegment does not map to a route, it tries using 'param' as a key.
// Path parameters are mapped with 'param' instead of the path segment (e.g. :id<int>).
func (r *route) getNestedRoute(pathSegment string) (*route, bool) {
	nestedRoute, ok := r.nestedRoutes[pathSegment]
	if !ok {
		nestedRoute, ok = r.nestedRoutes["param"]
	}
	return nestedRoute, ok
}

// makePathParam sets the path parameters value. isIntParam is set when '<int>'
// is provided. This tells the context object to convert the path paramter to
// an int. The name of the route is set to 'param' instead of the path segment.
func (r *route) makePathParam() {
	r.isPathParam = true
	if strings.HasSuffix(r.name, "<int>") {
		r.isIntParam = true
	}
	r.paramName = strings.TrimSuffix(r.name[1:], "<int>")
	r.name = "param"
}
