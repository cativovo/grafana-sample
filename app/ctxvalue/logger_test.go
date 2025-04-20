package ctxvalue_test

import (
	"app/ctxvalue"
	"context"
	"log/slog"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := ctxvalue.WithLogger(context.Background(), slog.Default())
	l := ctxvalue.Logger(ctx)
	if l == nil {
		t.Error("logger is nil")
	}
}
