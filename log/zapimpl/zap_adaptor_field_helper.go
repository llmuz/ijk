package zapimpl

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/llmuz/ijk/log"
)

type Option func(c *ZapHelperBuilder)

func AddHook(hook log.Hook) Option {
	return func(c *ZapHelperBuilder) {
		c.hooks.Add(hook)
	}
}

type ZapHelperBuilder struct {
	hooks log.LevelHooks
}

func NewHelper(logger *zap.Logger, opt ...Option) log.Helper {
	var cfg = &ZapHelperBuilder{hooks: make(log.LevelHooks, 6)}
	for _, v := range opt {
		v(cfg)
	}

	c := zapLoggerHelper{
		logger: logger,
		Hooks:  cfg.hooks,
	}
	c.initLogLevel()
	return &c
}

type zapLoggerHelper struct {
	logger *zap.Logger
	Hooks  log.LevelHooks
	level  log.Level
}

func (c *zapLoggerHelper) WithContext(ctx context.Context) log.FieldLogger {
	return &zapFieldLogger{
		ctx:    ctx,
		helper: c,
		entry:  NewZapLogEntry(ctx),
	}
}

func (c *zapLoggerHelper) levelEnabled(level log.Level) bool {
	return c.level <= level
}

// 初始化 level 值
func (c *zapLoggerHelper) initLogLevel() {
	var levels = []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
	}

	for _, v := range levels {
		if c.logger.Core().Enabled(v) {
			c.level = log.Level(v)
			break
		}
	}
}
