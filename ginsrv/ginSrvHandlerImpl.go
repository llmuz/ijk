package ginsrv

import (
	purger "errors"

	"github.com/gin-gonic/gin"

	"github.com/llmuz/ijk/errors"
)

func NewServiceHandler() ServiceHandler {
	return &ginSrvHandlerImpl{}
}

type ginSrvHandlerImpl struct {
}

func (c *ginSrvHandlerImpl) Error(ctx *gin.Context, err error) {
	var e errors.Error
	if purger.Is(err, &e) {
		c.response(ctx, e.GetHttpCode(), e.GetErrNo(), e.GetErrMsg(), e.GetData())
		return
	}

	_e := errors.New(500, errors.UnknownErrNo, "unknown", nil)
	c.response(ctx, _e.HttpCode, _e.GetErrNo(), _e.GetErrMsg(), _e.GetData())
}

func (c *ginSrvHandlerImpl) ParamsError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	e := errors.BadRequest("参数错误")
	c.response(ctx, e.GetHttpCode(), e.GetErrNo(), e.GetErrMsg(), e.GetData())
}

func (c *ginSrvHandlerImpl) Validate(ctx *gin.Context, req interface{}) (err error) {
	return err
}

func (c *ginSrvHandlerImpl) Success(ctx *gin.Context, data interface{}) {
	e := errors.Success(data)
	c.response(ctx, e.GetHttpCode(), e.GetErrNo(), e.GetErrMsg(), e.GetData())
}

func (c *ginSrvHandlerImpl) response(ctx *gin.Context, statusCode int32, errNo int64, errMsg string, data interface{}) {
	ctx.JSON(int(statusCode), errors.New(statusCode, errNo, errMsg, data))
}
