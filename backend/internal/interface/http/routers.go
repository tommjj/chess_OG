package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Group is a option function to group register router functions.
func Group(path string, registerRouterFuncs ...HTTPOptionFunc) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		g := r.Group(path)
		for _, fn := range registerRouterFuncs {
			fn(g)
		}
		return nil
	}
}

// RegisterPing registers a ping router for health check
func RegisterPing() HTTPOptionFunc {
	return func(r gin.IRouter) error {
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		return nil
	}
}

// RegisterStatic is a option function to return register static router function
func RegisterStatic(root string) HTTPOptionFunc {
	return func(i gin.IRouter) error {
		i.Static("/static", root)
		return nil
	}
}

// RegisterStaticFile is a option function to return register static file router function
func RegisterStaticFile(relativePath, filepath string) HTTPOptionFunc {
	return func(i gin.IRouter) error {
		i.StaticFile(relativePath, filepath)
		return nil
	}
}

// RegisterRoute registers a route with a gin.HandlerFunc
func RegisterRoute(path string, handler gin.HandlerFunc, methods ...string) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		switch len(methods) {
		case 0:
			r.GET(path, handler)
		default:
			for _, method := range methods {
				switch method {
				case "GET":
					r.GET(path, handler)
				case "POST":
					r.POST(path, handler)
				case "PUT":
					r.PUT(path, handler)
				case "DELETE":
					r.DELETE(path, handler)
				case "PATCH":
					r.PATCH(path, handler)
				case "HEAD":
					r.HEAD(path, handler)
				case "OPTIONS":
					r.OPTIONS(path, handler)
				}
			}
		}
		return nil
	}
}

// RegisterHTTPHandleRoute registers a route with a standard http.HandlerFunc
func RegisterHTTPHandleRoute(path string, handler http.HandlerFunc, methods ...string) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		switch len(methods) {
		case 0:
			r.GET(path, gin.WrapH(http.HandlerFunc(handler)))
		default:
			for _, method := range methods {
				switch method {
				case "GET":
					r.GET(path, gin.WrapH(http.HandlerFunc(handler)))
				case "POST":
					r.POST(path, gin.WrapH(http.HandlerFunc(handler)))
				case "PUT":
					r.PUT(path, gin.WrapH(http.HandlerFunc(handler)))
				case "DELETE":
					r.DELETE(path, gin.WrapH(http.HandlerFunc(handler)))
				case "PATCH":
					r.PATCH(path, gin.WrapH(http.HandlerFunc(handler)))
				case "HEAD":
					r.HEAD(path, gin.WrapH(http.HandlerFunc(handler)))
				case "OPTIONS":
					r.OPTIONS(path, gin.WrapH(http.HandlerFunc(handler)))
				}
			}
		}
		return nil
	}
}

func WithMiddleware(middleware ...gin.HandlerFunc) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		r.Use(middleware...)
		return nil
	}
}
