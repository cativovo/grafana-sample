package server

import (
	"app/ctxvalue"
	"app/otel"
	"context"
	"time"
)

func process2(ctx context.Context) {
	_, span := otel.Tracer.Start(ctx, "process2")
	defer span.End()

	logger := ctxvalue.Logger(ctx)

	logger.Info("Starting process 2")
	time.Sleep(time.Second * 1)
	logger.Info("Process 2 completed")
}

func process1(ctx context.Context) {
	ctx, span := otel.Tracer.Start(ctx, "process1")
	defer span.End()

	logger := ctxvalue.Logger(ctx)

	logger.Info("Starting process 1")
	time.Sleep(time.Second * 2)
	logger.Info("Process 1 completed")

	process2(ctx)
}
