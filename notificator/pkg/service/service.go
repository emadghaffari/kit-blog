package service

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
)

// NotificatorService describes the service.
type NotificatorService interface {
	// Add your methods here
	Send(ctx context.Context, to, body string) (id string, err error)
}

type basicNotificatorService struct{}

func (b *basicNotificatorService) Send(ctx context.Context, to string, body string) (id string, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			span := tracer.StartSpan("notification", opentracing.ChildOf(pctx))
			defer span.Finish()

			span.SetTag("id NOTIF", u.String())
		}
	}

	return u.String(), err
}

// NewBasicNotificatorService returns a naive, stateless implementation of NotificatorService.
func NewBasicNotificatorService() NotificatorService {
	return &basicNotificatorService{}
}

// New returns a NotificatorService with all of the expected middleware wired in.
func New(middleware []Middleware) NotificatorService {
	var svc NotificatorService = NewBasicNotificatorService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
