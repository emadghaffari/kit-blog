package service

import (
	"context"

	"github.com/emadghaffari/kit-blog/posts/model"
	"github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
)

// PostsService describes the service.
type PostsService interface {
	// Add your methods here
	Store(ctx context.Context, post model.Post) (response string, err error)
	Update(ctx context.Context, post model.Post) (response string, err error)
	List(ctx context.Context, post model.Post) (response []*pb.Post, err error)
	Delete(ctx context.Context, post model.Post) (response string, err error)
}

type basicPostsService struct{}

func (b *basicPostsService) Store(ctx context.Context, post model.Post) (response string, err error) {
	// TODO implement the business logic of Store
	return response, err
}
func (b *basicPostsService) Update(ctx context.Context, post model.Post) (response string, err error) {
	// TODO implement the business logic of Update
	return response, err
}
func (b *basicPostsService) List(ctx context.Context, post model.Post) (response []*pb.Post, err error) {
	// TODO implement the business logic of List
	return response, err
}
func (b *basicPostsService) Delete(ctx context.Context, post model.Post) (response string, err error) {
	// TODO implement the business logic of Delete
	return response, err
}

// NewBasicPostsService returns a naive, stateless implementation of PostsService.
func NewBasicPostsService() PostsService {
	return &basicPostsService{}
}

// New returns a PostsService with all of the expected middleware wired in.
func New(middleware []Middleware) PostsService {
	var svc PostsService = NewBasicPostsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
