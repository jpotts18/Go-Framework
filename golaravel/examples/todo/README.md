# Todo App Example

A simple todo application built with GoLaravel to test the framework's ergonomics.
Includes both an HTML frontend (server-rendered with templates) and a JSON API.

## Running

```sh
cd examples/todo
go run main.go
```

Visit http://localhost:5000 for the HTML app, or use the JSON API:

```sh
# List todos
curl http://localhost:5000/api/todos

# Create a todo
curl -X POST http://localhost:5000/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}'

# Toggle completion
curl -X PUT http://localhost:5000/api/todos/1/toggle

# Delete a todo
curl -X DELETE http://localhost:5000/api/todos/1
```

## Ergonomics Notes

Issues found and fixed while building this example:

### Bug: `Response.Redirect` panicked (fixed)

`Response.Redirect` passed `nil` as the `*http.Request` to `http.Redirect`, which
panics because `http.Redirect` reads the request method. Fixed by storing the
original `*http.Request` in the `Response` struct and passing it through.

### Bug: `req.All()` missed form POST data (fixed)

`Request.All()` only merged route params, query strings, and JSON bodies. It did
not parse `application/x-www-form-urlencoded` or `multipart/form-data` bodies,
so form submissions were invisible to validation. Fixed by adding form parsing.

### Observations

**What works well:**
- Route definition is clean and familiar (`router.Get`, `router.Post`, etc.)
- Route groups with middleware feel natural
- Validation API is expressive and the rules are easy to read
- Response chaining (`res.Status(201).JSON(data)`) is ergonomic
- Error responses (`res.NotFound`, `res.BadRequest`) are convenient
- Middleware is simple to apply both globally and per-route
- Handler signature `func(*Request, *Response) error` is clean

**Areas for improvement:**
- View rendering requires accessing `app.GetView().Render()` and then calling
  `res.HTML()` separately - a `res.View("name", data)` helper would be nicer,
  but would require solving the circular dependency between response and view
- No `res.Back()` or `res.Redirect().Back()` for redirecting to the previous
  page (common in form handling)
- `req.All()` returns `map[string]interface{}` - using this with validation
  works, but extracting typed values afterward is verbose (type assertions)
- Route parameter parsing returns strings only - a `req.ParamInt("id")` helper
  would reduce boilerplate `strconv.Atoi` calls
- No static file serving built into the router
