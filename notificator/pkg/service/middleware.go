package service

import (
	"context"
	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(NotificatorService) NotificatorService

type loggingMiddleware struct {
	logger log.Logger
	next   NotificatorService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a NotificatorService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next NotificatorService) NotificatorService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Send(ctx context.Context, to string, body string) (id string, err error) {
	defer func() {
		l.logger.Log("method", "Send", "to", to, "body", body, "id", id, "err", err)
	}()
	return l.next.Send(ctx, to, body)
}
