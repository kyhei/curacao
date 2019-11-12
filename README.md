# Curacao

Curacao is modular designed web framework for Golang.
Inspired by [Martini](https://github.com/go-martini/martini)

**Note: Curacao is heavy developing now. Don't use in production environment**

```go
package main

import "github.com/kyhei/curacao"

func main() {
  c := curacao.NewApp("0.0.0.0", "8080")
  c.Get(
    "/",
    func() string {
      return "Hello Curacao!!"
    },
  )
  
  c.Start()
}

```

## Routing

```go
c.Get("/", func(){}) // match http://localhost:8080

c.Get("/users", func(){}) // match http://localhost:8080/users

// match http://localhost:8080/users/123/show
c.Get("/users/:id/show", func(params curacao.HTTPParams){
  id, err := params.Get("id")
  if err != nil {
    fmt.Println("params not found")
  }
  fmt.Printf("received id is %v\n", id) // received id is 123
})
```

### HTTP Methods

```go
c.Get("/get", func(){}) // GET /get
c.POST("/post", func(){}) // POST /post
c.Register("DELETE", "/delete", func(){}) // DELETE /delete
```

`c.GET` and `c.Post` are short hand of `c.Register`

## Handlers

Handlers can receive any type of functions.  
Curacao can pass returned string or []byte to http.ResponseWriter.Write()  
And Curacao can write returned int to http.ResponseWriter.WriteHeader()  


```go
c.Get("/", func(){
  log.Println("Hello Curacao!")
}) // HTTP 200
```

```go
c.Get("/", func() string {
  return "Hello Curacao!"
}) // HTTP 200 Hello Curacao!
```

```go
c.Get("/", func() int, []byte {
  return 222, []byte("Hello Curacao!")
}) // HTTP 222 Hello Curacao!
```

### Service Injection

Because Curacao uses reflect to resolve handler,  
the arguments of handler are automatically injected by Curacao.  

```go
// GET /sample/caracao?lang=go
c.Get(
  "/sample/:name", 
  func(
    r *http.Request, 
    params curacao.HTTPParams, 
    query curacao.HTTPQuery,
  ){
    name, _ := params.Get("name")
    lang, _ := query.Get("lang")
    log.Printf("%v %v %v\n", r.URL.Path, name, lang)
    // -> /sample caracao go
  },
)
```

### Params

When you want to use params in handler,  
please define a function which has curacao.HTTPParams argument 

```go
c.Get(
  "/params/:hoge/:fuga", 
  func(params curacao.HTTPParams){
    hoge, err := params.Get("hoge")
    if err != nil {
      fmt.Println("not found")
    }

    fuga ,_ := params.Get("fuga")
    fmt.Printf("received params are %v %v\n", hoge, fuga)
  },
)
```

### Query string

```go

// GET /query?name=curacao&lang=go

c.Get(
  "/query",
  func(query curacao.HTTPQuery) {
    name, err := query.Get("name")
    if err != nil {
      fmt.Println("not found")
    }

    lang, _ := query.Get("lang")

    fmt.Printf("received queries are %v %v\n", name, lang)
  },
)

```

## Middleware

Middleware is called before handler.  
You can add middlewares easily.

```go
// Add Logger Middleware
logger := func(r *http.Request) {
  log.Printf("%v %v\n", r.Method, r.URL.Path)
}

c.Use(logger)
```

### Return values

Curacao determines whether to allow the request based on the return value from middleware

```go

// This middleware passes all request.
mw := func() (bool, int, string){
  ok := true // when you want to reject the request, return false
  status := 200 // http status code
  response := "middleware response" // middleware can return interface{}
  return ok, status, response
}
c.Use(mw)

// This middleware rejects all request
forbidden := func() (bool, int){
  return false, 403
}
c.Use(forbidden)

```

Below is example middlewara of Basic authentication.

```go
basicAuth := func(r *http.Request) (bool, int) {
  username, password, ok := r.BasicAuth()

  if ok && username == "curacao" && password == "password" {
    // accept request and call next middleware or handler
    return true, 200
  }

  return false, 403 // reject request
}

c.Use(basicAuth)
```

### Service Injection

Handler function can receive response from middleware.

```go

c.Use(func() string { 
  return "from middleware" 
})

type User struct {
  name string
}

c.Use(func() *User {
  user := new(User)
  user.name = "curacao"
  return user
})

c.Get(
  "/",
  func(s string, u *User){
    fmt.Printf("s is %v\n", s) // s is from middleware
    fmt.Printf("name of u is %v\n", u.name) // name of u is curacao
  }
)

```

## License
Curacao is distributed by The MIT License, see LICENSE