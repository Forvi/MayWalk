package log

import (
	"context"
	"log/slog"

	ctxUtil "github.com/Forvi/MayWalk/internal/common/ctx"
)

type ctxHandler struct {
	h slog.Handler
}

func (c ctxHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return c.h.Enabled(ctx, l)
}

func (c ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := ctxUtil.RequestID(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return c.h.Handle(ctx, r)
}

func (c ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ctxHandler{c.h.WithAttrs(attrs)}
}

func (c ctxHandler) WithGroup(name string) slog.Handler {
	return ctxHandler{c.h.WithGroup(name)}
}
