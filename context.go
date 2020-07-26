package gearbox

import (
	"github.com/valyala/fasthttp"
)

// Context interface
type Context interface {
	Next()
	Context() *fasthttp.RequestCtx
	Param(key string) string
	Query(key string) string
	SendString(value string) Context
	Status(status int) Context
	Set(key string, value string)
	Get(key string) string
	Body() string
}

// handlerFunc defines the handler used by middleware as return value.
type handlerFunc func(ctx Context)

// handlersChain defines a handlerFunc array.
type handlersChain []handlerFunc

// Context defines the current context of request and handlers/middlewares to execute
type context struct {
	requestCtx  *fasthttp.RequestCtx
	paramValues map[string]string
	handlers    handlersChain
	index       int
}

// Next function is used to successfully pass from current middleware to next middleware.
// if the middleware thinks it's okay to pass it
func (ctx *context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
	}
}

// Param returns value of path parameter specified by key
func (ctx *context) Param(key string) string {
	return ctx.paramValues[key]
}

// Context returns Fasthttp context
func (ctx *context) Context() *fasthttp.RequestCtx {
	return ctx.requestCtx
}

// SendString sets body of response as a string
func (ctx *context) SendString(value string) Context {
	ctx.requestCtx.SetBodyString(value)
	return ctx
}

// Status sets the HTTP status code
func (ctx *context) Status(status int) Context {
	ctx.requestCtx.Response.SetStatusCode(status)
	return ctx
}

// Get returns the HTTP request header specified by field key
func (ctx *context) Get(key string) string {
	return GetString(ctx.requestCtx.Request.Header.Peek(key))
}

// Set sets the response's HTTP header field key to the specified key, value
func (ctx *context) Set(key, value string) {
	ctx.requestCtx.Response.Header.Set(key, value)
}

// Query returns the query string parameter in the request url
func (ctx *context) Query(key string) string {
	return GetString(ctx.requestCtx.QueryArgs().Peek(key))
}

// Body returns the raw body submitted in a POST request
func (ctx *context) Body() string {
	return GetString(ctx.requestCtx.Request.Body())
}
