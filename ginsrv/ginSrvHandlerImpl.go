package ginsrv

import (
	"github.com/gin-gonic/gin"

	"github.com/llmuz/ijk/errors"
)

func NewServiceHandler(validate func(ctx *gin.Context, req interface{}) (err error)) ServiceHandler {
	return &ginSrvHandlerImpl{
		validate: validate,
	}
}

type ginSrvHandlerImpl struct {
	validate func(ctx *gin.Context, req interface{}) (err error)
}

func (c *ginSrvHandlerImpl) Error(ctx *gin.Context, err error) {
	if e, ok := err.(*errors.Error); ok {
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
	return c.validate(ctx, req)
}

func (c *ginSrvHandlerImpl) Success(ctx *gin.Context, data interface{}) {
	e := errors.Success(data)
	c.response(ctx, e.GetHttpCode(), e.GetErrNo(), e.GetErrMsg(), e.GetData())
}

func (c *ginSrvHandlerImpl) response(ctx *gin.Context, statusCode int32, errNo int64, errMsg string, data interface{}) {
	ctx.JSON(int(statusCode), errors.New(statusCode, errNo, errMsg, data))
}
