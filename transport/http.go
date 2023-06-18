package transport

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

func Middleware(m ...gin.HandlerFunc) ServerOption {
	return func(o *Server) {
		o.engine.Use(m...)
	}
}

func RunMode(runMode string) ServerOption {
	return func(o *Server) {
		gin.SetMode(runMode)
	}
}

func Endpoint(endpoint string) ServerOption {
	return func(o *Server) {
		o.srv.Addr = endpoint
	}
}

type Server struct {
	engine *gin.Engine  // gin 框架引擎
	srv    *http.Server // http 框架
}

func NewHttpServer(opts ...ServerOption) (srv *Server) {
	var engine = gin.New()
	s := &http.Server{
		Handler: engine,
	}

	srv = &Server{
		engine: engine,
		srv:    s,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (c *Server) Engine() *gin.Engine {
	return c.engine
}

func (c *Server) Start(ctx context.Context) (err error) {
	c.srv.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}

	err = c.srv.ListenAndServe()
	if err != nil {
		return err
	}

	return err
}

func (c *Server) Stop(ctx context.Context) (err error) {
	return c.srv.Shutdown(ctx)
}
