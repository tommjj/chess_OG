package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const EnvProduction = "production"

// Server
type Server interface {
	// Serve starts the server
	// blocking
	Serve() error

	// ServeTLS starts the server with TLS
	ServeTLS(certFile, keyFile string) error

	// Shutdown gracefully shuts down the server
	Shutdown(ctx context.Context) error

	// ShutdownWithTimeout gracefully shuts down the server with timeout
	ShutdownWithTimeout(d time.Duration) error
}

// HTTPOptionFunc defines a function type for configuring the HTTP router
type HTTPOptionFunc func(gin.IRouter) error

// Router is the HTTP server implementation
type Router struct {
	engine *gin.Engine
	server *http.Server
	Addr   string
}

// HTTPConfig holds configuration settings for the HTTP server
type HTTPConfig struct {
	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *tls.Config

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If zero, the value of
	// ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	ReadHeaderTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If zero, the value
	// of ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int

	// TLSNextProto optionally specifies a function to take over
	// ownership of the provided TLS connection when an ALPN
	// protocol upgrade has occurred. The map key is the protocol
	// name negotiated. The Handler argument should be used to
	// handle HTTP requests and will initialize the Request's TLS
	// and RemoteAddr if not already set. The connection is
	// automatically closed when the function returns.
	// If TLSNextProto is not nil, HTTP/2 support is not enabled
	// automatically.
	TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)

	// HTTP2 configures HTTP/2 connections.
	//
	// This field does not yet have any effect.
	// See https://go.dev/issue/67813.
	HTTP2 *http.HTTP2Config

	// Protocols is the set of protocols accepted by the server.
	//
	// If Protocols includes UnencryptedHTTP2, the server will accept
	// unencrypted HTTP/2 connections. The server can serve both
	// HTTP/1 and unencrypted HTTP/2 on the same address and port.
	//
	// If Protocols is nil, the default is usually HTTP/1 and HTTP/2.
	// If TLSNextProto is non-nil and does not contain an "h2" entry,
	// the default is HTTP/1 only.
	Protocols *http.Protocols
}

// NewAdapter creates a new HTTP server adapter
func NewAdapter(addr string, conf *HTTPConfig, options ...HTTPOptionFunc) (Server, error) {
	if os.Getenv("ENV") == EnvProduction { // gin.ReleaseMode
		gin.SetMode(gin.ReleaseMode)
	} else {
		fmt.Println("Running in development mode")
	}

	r := gin.New()

	for _, option := range options {
		err := option(r)
		if err != nil {
			return nil, fmt.Errorf("failed to apply HTTP option: %w", err)
		}
	}

	if conf == nil {
		conf = &HTTPConfig{}
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		ReadTimeout:       conf.ReadTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
		TLSConfig:         conf.TLSConfig,
		MaxHeaderBytes:    conf.MaxHeaderBytes,
		TLSNextProto:      conf.TLSNextProto,
		Protocols:         conf.Protocols,
		HTTP2:             conf.HTTP2,
	}

	return &Router{
		server: srv,
		engine: r,
		Addr:   addr,
	}, nil
}

// Serve starts the server
// blocking
func (r *Router) Serve() error {
	return r.server.ListenAndServe()
}

// ServeTLS starts the server with TLS
func (r *Router) ServeTLS(certFile, keyFile string) error {
	return r.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown gracefully shuts down the server
func (r *Router) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

// ShutdownWithTimeout gracefully shuts down the server with timeout
func (r *Router) ShutdownWithTimeout(d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return r.server.Shutdown(ctx)
}
