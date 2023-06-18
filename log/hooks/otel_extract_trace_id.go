package hooks

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/llmuz/ijk/log"
)

func NewOtelLogHook(lvs []log.Level) (h log.Hook) {
	return &otellogHook{lvs: lvs}
}

type otellogHook struct {
	lvs []log.Level
}

func (c *otellogHook) Levels() (lvs []log.Level) {
	return c.lvs
}

func (c *otellogHook) Fire(e log.Entry) (err error) {
	spanCtx := trace.SpanContextFromContext(e.Context())
	if !spanCtx.IsValid() {
		if err = e.AppendField(log.Any("trace_id", "")); err != nil {
			return err
		}

		if err = e.AppendField(log.Any("span_id", "")); err != nil {
			return err
		}
		return err
	}
	// 从 SpanContext 中提取 Trace ID 和 Span ID
	if err = e.AppendField(log.Any("trace_id", spanCtx.TraceID())); err != nil {
		return err
	}

	if err = e.AppendField(log.Any("span_id", spanCtx.SpanID())); err != nil {
		return err
	}
	return err
}
