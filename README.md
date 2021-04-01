# Go-Router

## Features
* Middleware
* Path Parameters
* Groups with different sets of middleware

## Installation

```bash
$ go get github.com/andromeda-voyager/go-router
```

## Usage

```go
import (
  "log"
  "net/http"

  router "github.com/andromeda-voyager/go-router"
)

func main() {

  http.HandleFunc("/", router.Handler)
  log.Fatal(http.ListenAndServe(":8080", nil))

  group := router.NewGroup()

  group.Post("/accounts", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
	
  }
  
}
```


### Middleware

Middleware has a similar method signature to a route callback with one exception, it returns a bool instead of void. Return true if the router should continue processing other middleware and eventually call the callback for that route. Return false if it should not continue (e.g., failed to authenticate the user). The router context can be used to store values and is accessible to other middleware and the route callback function. 
```go
func Auth(w http.ResponseWriter, r *http.Request, c *router.Context) bool {
  c.Keys["id"] = "1234"
  return true
}

func init() {
  authGroup := router.NewGroup()
  authGroup.Use(Auth)
  
  group.Post("/accounts", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
    id := c.Keys["id"] // retreive the id from the router context
  }
}
```
Different sets of middleware can be assigned to each router group that is created. 

### Path Parameters

Path parameters can be integers or strings. To let the router know there is a path parameter, prefix it with `:`. The router can also check if the paramter value provided is an integer if `<int>` is added as a suffix to the path paramter. When `<int>` is used and a call uses an invalid integer, the router will send a bad request response.
```go
  group.Post("/accounts/:id<int>", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
    accountID := c.Keys["id"].(int)
  }
  
  group.Post("/products/:productName", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
    productName := c.Keys["productName"]
  }
```
Routes with the same subpath can not use different path parameter names or types within those subpaths. The following routes should not be used together:
`"/user/:lastName"` and `"/user/:firstName"`. 
If the route with the `lastName` paramter is added to the router first, the router will not store the `firstName` value in the context with `firstName` as the key. However, it is retreievable using `lastName` as the key. For clarification, the following routes and their path paramters will work as expected: `"/user/:lastName"` and `"/user/:lastName/orders/:orderID"`. 
