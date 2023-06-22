package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/llmuz/ijk/errors"
	v1 "github.com/llmuz/ijk/examples/api/greeter/v1"
	"github.com/llmuz/ijk/ginsrv"
	"github.com/llmuz/ijk/log"
	"github.com/llmuz/ijk/log/hooks"
	"github.com/llmuz/ijk/log/zapimpl"
	"github.com/llmuz/ijk/middleware/logging"
	"github.com/llmuz/ijk/middleware/ratelimit"
	"github.com/llmuz/ijk/middleware/recovery"
	"github.com/llmuz/ijk/middleware/tracing"
	"github.com/llmuz/ijk/transport"
)

type ValidateHandler struct {
	log      log.Helper
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidateHandler(log log.Helper) (v *ValidateHandler) {
	var validate = validator.New()
	var trans, found = ut.New(zh.New()).GetTranslator("zh")
	if !found {
		log.WithContext(context.TODO()).Infof("not found translator zh")
	}
	var err = zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.WithContext(context.TODO()).Fatalf("register translator error")
	}

	v = &ValidateHandler{validate: validate, trans: trans, log: log}

	return v
}

func (c *ValidateHandler) Validate(ctx *gin.Context, req interface{}) (err error) {
	if err = c.validate.Struct(req); err != nil {
		c.log.WithContext(ctx.Request.Context()).Infof("req %#v", req)
		if transErr, ok := err.(validator.ValidationErrors); ok {

			translations := transErr.Translate(c.trans)
			var buff = make([]string, 0, len(translations))

			for k, v := range translations {
				buff = append(buff, fmt.Sprintf("%s = %s", k, v))
			}
			request := errors.New(400, 400, "参数错误", buff)
			c.log.WithContext(ctx.Request.Context()).Infof("error %s, buff %v  %v", translations, buff, request)
			return request

		} else {
			c.log.WithContext(ctx.Request.Context()).Infof("error %s", err)
			return errors.BadRequest("bad params")

		}

	}
	return err
}

type greeterImpl struct {
	log log.Helper
}

func NewGreeterServer(log log.Helper) v1.GreeterGin {
	return &greeterImpl{log: log}
}

func (c *greeterImpl) Greeter(ctx context.Context, req *v1.GreeterRequest) (resp *v1.GreeterResponse, err error) {
	c.log.WithContext(ctx).Infof("req %s", req.String())
	return &v1.GreeterResponse{Data: req.GetName()}, nil
}

func main() {

	logger, err := zapimpl.Logger(
		zapimpl.Compress(true),
		zapimpl.MaxSize(1024),
		zapimpl.MaxBackup(10),
		zapimpl.FileName("biz.log"),
		zapimpl.Level("debug"),
		zapimpl.MaxAge(3),
		zapimpl.JsonFormat(true),
		zapimpl.DebugModeOutputConsole(true),
		zapimpl.EncoderConfig(zapimpl.DefaultEncoderConfig),
	)
	if err != nil {
		panic(err)
	}

	// http 服务
	srv := transport.NewHttpServer(
		transport.RunMode(gin.DebugMode),
		transport.Endpoint("127.0.0.1:8081"),
		transport.Middleware(
			gin.CustomRecovery(
				recovery.ServerPanic(
					zapimpl.NewHelper(
						logger.WithOptions(zap.AddCallerSkip(2)),
						zapimpl.AddHook(hooks.NewOtelLogHook(log.DefaultLevel)),
					),
				),
			),
			ratelimit.Server(),
			tracing.WithTraceProvider(
				tracesdk.NewTracerProvider(
					tracesdk.WithSampler(tracesdk.NeverSample()),
				),
			),
			logging.ServerLog(
				zapimpl.NewHelper(
					logger.WithOptions(zap.AddCallerSkip(2)),
					zapimpl.AddHook(
						hooks.NewOtelLogHook(log.DefaultLevel),
					),
				),
			),
		),
	)

	greeterServer := NewGreeterServer(
		zapimpl.NewHelper(
			logger.WithOptions(zap.AddCallerSkip(2)),
			zapimpl.AddHook(hooks.NewOtelLogHook(log.DefaultLevel)),
		),
	)

	v := NewValidateHandler(
		zapimpl.NewHelper(
			logger.WithOptions(zap.AddCallerSkip(2)),
			zapimpl.AddHook(hooks.NewOtelLogHook(log.DefaultLevel)),
		),
	)
	greeterSrvImpl := v1.NewGreeterHandler(greeterServer, ginsrv.NewServiceHandler(v.Validate))
	srv.Engine().GET("/greeter/v1/:name", greeterSrvImpl.GreeterHandler)

	srv.Engine().GET("/hello", func(ctx *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Int63n(100)))
		ctx.JSON(200, errors.Success("ok"))
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
