package repository

import (
	"app"
	"app/ctxvalue"
	"app/service"
	"context"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	trace trace.Tracer

	dbMu sync.RWMutex
	db   map[int]service.Something

	idMU sync.Mutex
	id   int
}

var _ service.Repository = (*Repository)(nil)

func NewRepository(t trace.Tracer) *Repository {
	return &Repository{
		db:    make(map[int]service.Something),
		trace: t,
	}
}

func (r *Repository) GetSomething(ctx context.Context, id string) (service.Something, error) {
	_, span := r.trace.Start(ctx, "Repository.GetSomething")
	defer span.End()

	logger := ctxvalue.Logger(ctx)

	logger.Info("Getting Something from db", "id", id)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return service.Something{}, app.NewError(app.ErrCodeNotFound, "Something not found")
	}

	time.Sleep(time.Second)

	r.dbMu.RLock()
	defer r.dbMu.RUnlock()

	s, ok := r.db[idInt]
	if !ok {
		return service.Something{}, app.NewError(app.ErrCodeNotFound, "Something not found")
	}

	return s, nil
}

func (r *Repository) CreateSomething(ctx context.Context, c service.CreateSomethingReq) (service.Something, error) {
	_, span := r.trace.Start(ctx, "Repository.CreateSomething")
	defer span.End()

	logger := ctxvalue.Logger(ctx)
	logger.Info("Creating Something into db")

	time.Sleep(time.Second * 2)

	r.idMU.Lock()
	defer r.idMU.Unlock()

	id := r.id
	r.id++

	s := service.Something{
		ID:    strconv.Itoa(id),
		Value: c.Value,
	}

	r.dbMu.Lock()
	defer r.dbMu.Unlock()

	r.db[id] = s

	return s, nil
}
