package http

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type HTTPOptionFunc func(gin.IRouter) error

type router struct {
	engine *gin.Engine
	Addr   string
}

// NewAdapter creates a new HTTP server adapter
func NewAdapter(addr string, options ...HTTPOptionFunc) (*router, error) {
	if os.Getenv("ENV") == "production" { // gin.ReleaseMode
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Swagger
	// docs.SwaggerInfo.BasePath = "/v1/api"
	// r.GET("/v1/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, option := range options {
		err := option(r)
		if err != nil {
			return nil, fmt.Errorf("failed to apply HTTP option: %w", err)
		}
	}

	return &router{
		engine: r,
		Addr:   addr,
	}, nil
}

// Serve is a method start server
//
// blocking
func (r *router) Serve() error {
	fmt.Printf("Starting HTTP server at %s\n", r.Addr)

	if err := r.engine.Run(r.Addr); err != nil {
		return err
	}
	return nil
}
