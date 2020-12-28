package endpoint

import (
	"context"

	endpoint "github.com/go-kit/kit/endpoint"

	"github.com/emadghaffari/kit-blog/comments/pkg/grpc/pb"
	service "github.com/emadghaffari/kit-blog/comments/pkg/service"
)

// StoreRequest collects the request parameters for the Store method.
type StoreRequest struct {
	Cm service.Comment `json:"cm"`
}

// StoreResponse collects the response parameters for the Store method.
type StoreResponse struct {
	Id  string `json:"id"`
	Err error  `json:"err"`
}

// MakeStoreEndpoint returns an endpoint that invokes Store on the service.
func MakeStoreEndpoint(s service.CommentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreRequest)
		id, err := s.Store(ctx, req.Cm)
		return StoreResponse{
			Err: err,
			Id:  id,
		}, nil
	}
}

// Failed implements Failer.
func (r StoreResponse) Failed() error {
	return r.Err
}

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Cm service.Comment `json:"cm"`
}

// UpdateResponse collects the response parameters for the Update method.
type UpdateResponse struct {
	Id  string `json:"id"`
	Err error  `json:"err"`
}

// MakeUpdateEndpoint returns an endpoint that invokes Update on the service.
func MakeUpdateEndpoint(s service.CommentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		id, err := s.Update(ctx, req.Cm)
		return UpdateResponse{
			Err: err,
			Id:  id,
		}, nil
	}
}

// Failed implements Failer.
func (r UpdateResponse) Failed() error {
	return r.Err
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
	PostID string `json:"post_id"`
}

// ListResponse collects the response parameters for the List method.
type ListResponse struct {
	CMS []*pb.Comment `json:"cms"`
	Err error         `json:"err"`
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
func MakeListEndpoint(s service.CommentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		cms, err := s.List(ctx, req.PostID)
		return ListResponse{
			CMS: cms,
			Err: err,
		}, nil
	}
}

// Failed implements Failer.
func (r ListResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Store implements Service. Primarily useful in a client.
func (e Endpoints) Store(ctx context.Context, cm service.Comment) (id string, err error) {
	request := StoreRequest{Cm: cm}
	response, err := e.StoreEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(StoreResponse).Id, response.(StoreResponse).Err
}

// Update implements Service. Primarily useful in a client.
func (e Endpoints) Update(ctx context.Context, cm service.Comment) (id string, err error) {
	request := UpdateRequest{Cm: cm}
	response, err := e.UpdateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(UpdateResponse).Id, response.(UpdateResponse).Err
}

// List implements Service. Primarily useful in a client.
func (e Endpoints) List(ctx context.Context, postID string) (CMS []*pb.Comment, err error) {
	request := ListRequest{PostID: postID}
	response, err := e.ListEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(ListResponse).CMS, response.(ListResponse).Err
}
