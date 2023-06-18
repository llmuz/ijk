package main

import (
	"context"
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/llmuz/ijk/transport"
)

func main() {
	// http 服务
	srv := transport.NewHttpServer(
		transport.RunMode(gin.DebugMode),
		transport.Endpoint("127.0.0.1:8081"),
		transport.Middleware(
			gin.Recovery(),
			gin.Logger(),
		),
	)

	srv.Engine().GET("/hello", func(ctx *gin.Context) {

	})

	srv.Engine().GET("/hello/2", func(ctx *gin.Context) {

	})
	pprof.Register(srv.Engine())
	if err := srv.Start(context.TODO()); err != nil {
		fmt.Println(err)
	}

	defer srv.Stop(context.TODO())
}
