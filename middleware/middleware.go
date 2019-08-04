package middleware

import (
	"net/http"
)

//HTTPMiddleware for func(ResponseWriter, *Request)
type HTTPMiddleware func(h http.HandlerFunc) http.HandlerFunc

//Pipeline chain middlewares together
type Pipeline interface {
	Do(a HTTPMiddleware) Pipeline
	For(h http.HandlerFunc) http.HandlerFunc
}

type pipe struct {
	middlewares []HTTPMiddleware
}

//NewPipeline build a middleware chain
func NewPipeline() Pipeline {
	return &pipe{middlewares: []HTTPMiddleware{}}
}

func (c *pipe) Do(a HTTPMiddleware) Pipeline {
	c.middlewares = append([]HTTPMiddleware{a}, c.middlewares...)
	return c
}

func (c *pipe) For(h http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range c.middlewares {
		h = middleware(h)
	}

	return h
}
