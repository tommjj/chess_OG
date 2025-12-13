package http

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// WithLogger adds logging middleware to the HTTP router
// ! Currently not used
func WithLogger() HTTPOptionFunc {
	return func(r gin.IRouter) error {
		// // set logger middleware
		// logger, err := logx.New(conf)
		// if err != nil {
		// 	panic(err)
		// }

		// r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		// r.Use(ginzap.RecoveryWithZap(logger, true))
		return nil
	}
}

// WithRecovery adds recovery middleware to the HTTP router
func WithRecovery() HTTPOptionFunc {
	return func(r gin.IRouter) error {
		r.Use(gin.Recovery())
		return nil
	}
}

// WithCORS adds CORS middleware to the HTTP router
//
//	allowedOrigins: list of allowed origins
//	allowHeaders: list of allowed headers
//	allowCredentials: whether to allow credentials
func WithCORS(allowedOrigins []string, allowHeaders []string, allowCredentials bool) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		CORSConfig := cors.DefaultConfig()

		CORSConfig.AllowCredentials = allowCredentials
		if allowCredentials { // when credentials are allowed, specific origins must be set
			CORSConfig.AllowOrigins = allowedOrigins
		} else if len(allowedOrigins) == 0 { // when no origins are specified and credentials are not allowed, allow all origins
			CORSConfig.AllowAllOrigins = true
		}

		for _, h := range allowHeaders {
			CORSConfig.AddAllowHeaders(h)
		}
		r.Use(cors.New(CORSConfig))
		return nil
	}
}

// WithNoCache adds no-cache middleware to the HTTP router
func WithNoCache() HTTPOptionFunc {
	return func(r gin.IRouter) error {
		r.Use(func(c *gin.Context) {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.Header("Surrogate-Control", "no-store")
			c.Next()
		})
		return nil
	}
}

// WithCustomValidators adds custom validators to the HTTP router
func WithCustomValidators(tag string, fn validator.Func, callValidationEvenIfNull ...bool) HTTPOptionFunc {
	return func(r gin.IRouter) error {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return fmt.Errorf("validator engine is not of type *validator.Validate")
		}

		err := v.RegisterValidation(tag, fn, callValidationEvenIfNull...)
		if err != nil {
			return fmt.Errorf("failed to register custom validator with tag %s: %w", tag, err)
		}
		return nil
	}
}
