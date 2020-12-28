package service

import (
	"context"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(UsersService) UsersService

type loggingMiddleware struct {
	logger log.Logger
	next   UsersService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UsersService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UsersService) UsersService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Login(ctx context.Context, username string, password string) (s0 string, e1 error) {
	defer func() {
		l.logger.Log("method", "Login", "username", username, "password", password, "s0", s0, "e1", e1)
	}()
	return l.next.Login(ctx, username, password)
}
func (l loggingMiddleware) Register(ctx context.Context, username string, password string, email string, phone string) (s0 string, e1 error) {
	defer func() {
		l.logger.Log("method", "Register", "username", username, "password", password, "email", email, "phone", phone, "s0", s0, "e1", e1)
	}()
	return l.next.Register(ctx, username, password, email, phone)
}

func (l loggingMiddleware) Get(ctx context.Context, id string) (s0, s1, s2 string, e1 error) {
	defer func() {
		l.logger.Log("method", "Get", "id", id, "s0", s0, "s1", s1, "s2", s2, "e1", e1)
	}()
	return l.next.Get(ctx, id)
}
