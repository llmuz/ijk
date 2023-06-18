package ginsrv

import "github.com/gin-gonic/gin"

type ServiceHandler interface {
	Error(ctx *gin.Context, err error)
	ParamsError(ctx *gin.Context, err error)
	Validate(ctx *gin.Context, req interface{}) (err error)
	Success(ctx *gin.Context, data interface{})
}
