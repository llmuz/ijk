// Code generated by github.com/llmuz/ijk/cmd/protoc-gen-go-gin. DO NOT EDIT.

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	ginsrv "github.com/llmuz/ijk/ginsrv"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the github.com/llmuz/ijk/cmd/protoc-gen-go-gin package it is being compiled against.
// context.metadata.
//gin.
//ginsrv.

// generate by protoc-gen-go-gin
// powered by xx

type GreeterGinServer interface {

	// /greeter/v1/:name GET
	Greeter(ctx context.Context, req *GreeterRequest) (resp *GreeterResponse, err error)
}

func NewGreeterServer(srv GreeterGinServer, srvHandler ginsrv.ServiceHandler) (c *Greeter) {
	c = &Greeter{
		server:     srv,
		srvHandler: srvHandler,
	}
	return c
}

type Greeter struct {
	server     GreeterGinServer
	srvHandler ginsrv.ServiceHandler
}

func (s *Greeter) Greeter_0(ctx *gin.Context) {
	var in GreeterRequest

	if err := ctx.ShouldBindUri(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&in); err != nil {
		s.srvHandler.ParamsError(ctx, err)
		return
	}

	if err := s.srvHandler.Validate(ctx, &in); err != nil {
		s.srvHandler.Error(ctx, err)
		return
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_context", ctx))

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterGinServer).Greeter(newCtx, &in)
	if err != nil {
		s.srvHandler.Error(ctx, err)
		return
	}

	s.srvHandler.Success(ctx, out)
}
