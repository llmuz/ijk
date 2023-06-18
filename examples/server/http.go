package main

import (
	"context"
	"fmt"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/llmuz/ijk/log"
	"github.com/llmuz/ijk/log/hooks"
	"github.com/llmuz/ijk/log/zapimpl"
	"github.com/llmuz/ijk/middleware/logging"
	"github.com/llmuz/ijk/middleware/recovery"
	"github.com/llmuz/ijk/middleware/tracing"
	"github.com/llmuz/ijk/transport"
)

func main() {
	logger, _ := zap.NewProduction()
	helper := zapimpl.NewHelper(
		logger,
		zapimpl.AddHook(hooks.NewOtelLogHook([]log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.WarnLevel,
			log.ErrorLevel,
			log.PanicLevel,
		}),
		),
	)
	// http 服务
	srv := transport.NewHttpServer(
		transport.RunMode(gin.DebugMode),
		transport.Endpoint("127.0.0.1:8081"),
		transport.Middleware(
			gin.CustomRecovery(recovery.ServerPanic(helper)),
			tracing.WithTraceProvider(tracesdk.NewTracerProvider(tracesdk.WithSampler(tracesdk.NeverSample()))),
			logging.ServerLog(helper),
		),
	)

	srv.Engine().GET("/hello", func(ctx *gin.Context) {

	})

	srv.Engine().GET("/panic", func(ctx *gin.Context) {
		panic("panic error")
	})
	pprof.Register(srv.Engine())
	if err := srv.Start(context.TODO()); err != nil {
		fmt.Println(err)
	}

	defer srv.Stop(context.TODO())
}
