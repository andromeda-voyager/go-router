package router

import "strconv"

// Context is used to store values set by the middleware and path parameters
type Context struct {
	Keys map[string]interface{}
}

// addPathParam sets a context value for a path parameter used in a route. If
// the path paramter is supposed to be an int, it is converted to an int.
func (c *Context) addPathParam(r *route, value string) bool {
	if r.isIntParam {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		c.Keys[r.paramName] = intValue
	} else {
		c.Keys[r.paramName] = value
	}
	return true
}
