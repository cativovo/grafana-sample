package service

import (
	"app"
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Something struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type CreateSomethingReq struct {
	Value string `json:"value"`
}

type Service interface {
	GetSomething(ctx context.Context, id string) (Something, error)
	CreateSomething(ctx context.Context, c CreateSomethingReq) (Something, error)
}

type Repository interface {
	GetSomething(ctx context.Context, id string) (Something, error)
	CreateSomething(ctx context.Context, c CreateSomethingReq) (Something, error)
}

type service struct {
	tracer trace.Tracer
	repo   Repository
}

func NewService(t trace.Tracer, r Repository) Service {
	return &service{
		tracer: t,
		repo:   r,
	}
}

func (s *service) GetSomething(ctx context.Context, id string) (Something, error) {
	_, span := s.tracer.Start(ctx, "service.GetSomething")
	defer span.End()

	time.Sleep(time.Millisecond * 1500)

	something, err := s.repo.GetSomething(ctx, id)
	if err != nil {
		return Something{}, fmt.Errorf("service.GetSomething: %w", err)
	}

	return something, nil
}

func (s *service) CreateSomething(ctx context.Context, c CreateSomethingReq) (Something, error) {
	_, span := s.tracer.Start(ctx, "service.GetSomething")
	defer span.End()

	if c.Value == "" {
		return Something{}, app.NewError(app.ErrCodeInvalid, "Value is empty")
	}

	something, err := s.repo.CreateSomething(ctx, c)
	if err != nil {
		return Something{}, fmt.Errorf("service.CreateSomething: %w", err)
	}

	return something, nil
}
