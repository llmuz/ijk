// Package ratelimit 备注, 代码来自于 github.com/go-kratos
package ratelimit

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	"github.com/llmuz/ijk/errors"
)

// ErrLimitExceed is service unavailable due to rate limit exceeded.
var ErrLimitExceed = errors.New(429, 429, "RATELIMIT service unavailable due to rate limit exceeded", nil)

// Option is ratelimit option.
type Option func(*options)

// WithLimiter set Limiter implementation,
// default is bbr limiter
func WithLimiter(limiter ratelimit.Limiter) Option {
	return func(o *options) {
		o.limiter = limiter
	}
}

type options struct {
	limiter ratelimit.Limiter
}

// Server ratelimiter middleware
func Server(opts ...Option) gin.HandlerFunc {
	o := &options{
		limiter: bbr.NewLimiter(),
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(ctx *gin.Context) {
		done, err := o.limiter.Allow()
		if err != nil {
			ctx.AbortWithStatusJSON(int(ErrLimitExceed.GetHttpCode()), ErrLimitExceed)
			return
		}
		ctx.Next()
		// allowed
		done(ratelimit.DoneInfo{Err: err})
		return
	}
}
