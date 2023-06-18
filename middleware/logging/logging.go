package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/llmuz/ijk/log"
)

func ServerLog(logger log.Helper) gin.HandlerFunc {
	// 初始化
	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			//isTerm:  isTerm,
			Keys: c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		sb := strings.Builder{}
		for _, v := range append(defaultLogFormatter(param)) {
			sb.WriteString(fmt.Sprintf("%s=%+v ", v.Key, v.Interface))
		}
		logger.WithContext(c.Request.Context()).Infof(sb.String())
	}
}

var defaultLogFormatter = func(param gin.LogFormatterParams) (fields []log.Field) {

	fields = make([]log.Field, 0, 8)
	fields = append(fields, log.Any("status_code", param.StatusCode))
	fields = append(fields, log.Any("latency_seconds", param.Latency.Seconds()))
	fields = append(fields, log.Any("client_ip", param.ClientIP))
	fields = append(fields, log.Any("method", param.Method))
	fields = append(fields, log.Any("path", param.Path))
	fields = append(fields, log.Any("body_size", param.BodySize))
	fields = append(fields, log.Any("error_message", param.ErrorMessage))
	return fields

}
