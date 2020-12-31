package service

import (
	"context"

	log "github.com/go-kit/kit/log"

	model "github.com/emadghaffari/kit-blog/posts/model"
	"github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
)

// Middleware describes a service middleware.
type Middleware func(PostsService) PostsService

type loggingMiddleware struct {
	logger log.Logger
	next   PostsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a PostsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next PostsService) PostsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Store(ctx context.Context, post model.Post) (response string, err error) {
	defer func() {
		l.logger.Log("method", "Store", "post", post, "response", response, "err", err)
	}()
	return l.next.Store(ctx, post)
}
func (l loggingMiddleware) Update(ctx context.Context, post model.Post) (response string, err error) {
	defer func() {
		l.logger.Log("method", "Update", "post", post, "response", response, "err", err)
	}()
	return l.next.Update(ctx, post)
}
func (l loggingMiddleware) List(ctx context.Context, post model.Post) (response []*pb.Post, err error) {
	defer func() {
		l.logger.Log("method", "List", "post", post, "response", response, "err", err)
	}()
	return l.next.List(ctx, post)
}
func (l loggingMiddleware) Delete(ctx context.Context, post model.Post) (response string, err error) {
	defer func() {
		l.logger.Log("method", "Delete", "post", post, "response", response, "err", err)
	}()
	return l.next.Delete(ctx, post)
}
