package service

import (
	"context"

	log "github.com/go-kit/kit/log"

	"github.com/emadghaffari/kit-blog/comments/pkg/grpc/pb"
)

// Middleware describes a service middleware.
type Middleware func(CommentsService) CommentsService

type loggingMiddleware struct {
	logger log.Logger
	next   CommentsService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a CommentsService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next CommentsService) CommentsService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Store(ctx context.Context, cm Comment) (id string, err error) {
	defer func() {
		l.logger.Log("method", "Store", "cm", cm, "id", id, "err", err)
	}()
	return l.next.Store(ctx, cm)
}
func (l loggingMiddleware) Update(ctx context.Context, cm Comment) (id string, err error) {
	defer func() {
		l.logger.Log("method", "Update", "cm", cm, "id", id, "err", err)
	}()
	return l.next.Update(ctx, cm)
}
func (l loggingMiddleware) List(ctx context.Context, postID string) (cms []*pb.Comment, err error) {
	defer func() {
		l.logger.Log("method", "List", "postID", postID, "cms", cms, "err", err)
	}()
	return l.next.List(ctx, postID)
}
