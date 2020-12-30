package service

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"

	"github.com/emadghaffari/kit-blog/comments/pkg/grpc/pb"
)

// Comment struct
type Comment struct {
	ID     string `json:"id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	PostID string `json:"post_id,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
}

// CommentsService describes the service.
type CommentsService interface {
	// Add your methods here
	Store(ctx context.Context, cm Comment) (id string, err error)
	Update(ctx context.Context, cm Comment) (id string, err error)
	List(ctx context.Context, postID string) (cms []*pb.Comment, err error)
}

type basicCommentsService struct{}

func (b *basicCommentsService) Store(ctx context.Context, cm Comment) (id string, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("store")

	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}

	span.Finish()
	return u.String(), err
}
func (b *basicCommentsService) Update(ctx context.Context, cm Comment) (id string, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}
	return u.String(), err
}
func (b *basicCommentsService) List(ctx context.Context, postID string) (cms []*pb.Comment, err error) {
	return []*pb.Comment{
		{
			PostID: "1",
			UserID: "2",
			Title:  "T",
			Body:   "B",
		},
		{
			PostID: "1",
			UserID: "2",
			Title:  "T",
			Body:   "B",
		},
	}, err
}

// NewBasicCommentsService returns a naive, stateless implementation of CommentsService.
func NewBasicCommentsService() CommentsService {
	return &basicCommentsService{}
}

// New returns a CommentsService with all of the expected middleware wired in.
func New(middleware []Middleware) CommentsService {
	var svc CommentsService = NewBasicCommentsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
