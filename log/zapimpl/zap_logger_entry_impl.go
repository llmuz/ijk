package zapimpl

import (
	"context"

	"github.com/llmuz/ijk/log"
)

type ZapLoggerEntry struct {
	ctx    context.Context // context
	fields []log.Field     // data gen by Entry impl
}

// Context get ctx
func (c *ZapLoggerEntry) Context() (ctx context.Context) {
	return c.ctx
}

// AppendField append field to fields
func (c *ZapLoggerEntry) AppendField(field log.Field) (err error) {
	c.fields = append(c.fields, field)
	return nil
}

// GetFields get fields
func (c *ZapLoggerEntry) GetFields() []log.Field {
	return c.fields
}

func NewZapLogEntry(ctx context.Context) log.Entry {
	return &ZapLoggerEntry{ctx: ctx, fields: make([]log.Field, 0, 4)}
}
