# HTTP Layer Refactor

## Why remove Gin?

The Smart Reconciliation Platform has been refactored to use Go's standard library (`net/http`) with `chi` instead of Gin.

Reasons:

- Fewer external dependencies
- Smaller API surface
- Better compatibility with Go tooling
- Easier testing
- Reduced framework lock-in
- More idiomatic Go

---

# Previous Design

Handlers depended directly on Gin.

```go
func (h *Handler) Login(c *gin.Context)
```

Authentication middleware also depended on Gin.

```go
func JWTAuth(...) gin.HandlerFunc
```

Business logic was therefore tightly coupled to a web framework.

---

# New Design

Handlers now use the standard Go HTTP handler.

```go
func (h *Handler) Login(
    w http.ResponseWriter,
    r *http.Request,
)
```

Every router in Go understands this signature.

Examples:

- net/http
- chi
- gorilla/mux
- httprouter

This makes the application independent of any specific framework.

---

# request.Decode()

Previously every handler contained:

```go
json.NewDecoder(r.Body).Decode(...)
```

The same code appeared repeatedly.

The Decode helper centralizes request parsing.

Benefits:

- Less duplication
- Centralized validation
- Easier future enhancements

Example:

```go
if err := request.Decode(r, &req); err != nil {
    response.BadRequest(w, err.Error())
    return
}
```

---

# decoder.DisallowUnknownFields()

This was intentionally enabled.

Without it:

```json
{
    "email":"john@example.com",
    "password":"secret",
    "abcdef":123
}
```

would silently ignore

```
abcdef
```

With DisallowUnknownFields enabled,

the request is rejected.

This prevents clients from accidentally sending invalid payloads.

---

# response.JSON()

Previously every handler repeated:

```go
w.Header().Set(...)
w.WriteHeader(...)
json.NewEncoder(...).Encode(...)
```

Now every response goes through one helper.

Benefits:

- Consistent responses
- Less boilerplate
- Easier maintenance

---

# Success helpers

Instead of:

```go
JSON(w, http.StatusCreated, ...)
```

handlers simply call

```go
response.Created(...)
```

Other helpers:

- OK
- Created
- Accepted
- NoContent

This improves readability.

---

# Error helpers

Instead of repeatedly writing

```go
map[string]string{
    "error": "...",
}
```

responses now use

```go
response.BadRequest(...)
response.Unauthorized(...)
response.Forbidden(...)
response.NotFound(...)
response.InternalServerError(...)
```

Advantages:

- Consistent JSON
- Centralized formatting
- Easier future changes

---

# Authentication Context

Previously

```go
c.Set("user_id", ...)
```

was used.

Now

```go
context.WithValue(...)
```

stores authentication information.

Handlers retrieve it using

```go
user, ok := authmiddleware.GetUser(r.Context())
```

This removes all Gin dependencies.

---

# Typed Context Keys

String keys were replaced with custom types.

Instead of

```go
"user_id"
```

we now use

```go
type contextKey string
```

Benefits:

- Prevents key collisions
- Recommended by the Go team
- Safer for large applications

---

# Typed Responses

Anonymous maps like

```go
map[string]any{
    ...
}
```

have been replaced by structs.

Example:

```go
type LoginResponse struct {
    Token string `json:"token"`
}
```

Benefits:

- Compile-time safety
- Better documentation
- Easier refactoring
- Cleaner API contracts

---

# Result

The HTTP layer now has clear responsibilities.

request/

Responsible for:

- Parsing requests

response/

Responsible for:

- Formatting responses

middleware/

Responsible for:

- Authentication
- Authorization
- Request lifecycle

handlers/

Responsible only for business logic.

This separation follows the Single Responsibility Principle and makes the codebase easier to test, extend, and maintain.