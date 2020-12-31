package endpoint

import (
	"context"

	endpoint "github.com/go-kit/kit/endpoint"

	model "github.com/emadghaffari/kit-blog/posts/model"
	"github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
	service "github.com/emadghaffari/kit-blog/posts/pkg/service"
)

// StoreRequest collects the request parameters for the Store method.
type StoreRequest struct {
	Post model.Post `json:"post"`
}

// StoreResponse collects the response parameters for the Store method.
type StoreResponse struct {
	Response string `json:"response"`
	Err      error  `json:"err"`
}

// MakeStoreEndpoint returns an endpoint that invokes Store on the service.
func MakeStoreEndpoint(s service.PostsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreRequest)
		response, err := s.Store(ctx, req.Post)
		return StoreResponse{
			Err:      err,
			Response: response,
		}, nil
	}
}

// Failed implements Failer.
func (r StoreResponse) Failed() error {
	return r.Err
}

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Post model.Post `json:"post"`
}

// UpdateResponse collects the response parameters for the Update method.
type UpdateResponse struct {
	Response string `json:"response"`
	Err      error  `json:"err"`
}

// MakeUpdateEndpoint returns an endpoint that invokes Update on the service.
func MakeUpdateEndpoint(s service.PostsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		response, err := s.Update(ctx, req.Post)
		return UpdateResponse{
			Err:      err,
			Response: response,
		}, nil
	}
}

// Failed implements Failer.
func (r UpdateResponse) Failed() error {
	return r.Err
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
	Post model.Post `json:"post"`
}

// ListResponse collects the response parameters for the List method.
type ListResponse struct {
	Response []*pb.Post `json:"response"`
	Err      error      `json:"err"`
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
func MakeListEndpoint(s service.PostsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		response, err := s.List(ctx, req.Post)
		return ListResponse{
			Err:      err,
			Response: response,
		}, nil
	}
}

// Failed implements Failer.
func (r ListResponse) Failed() error {
	return r.Err
}

// DeleteRequest collects the request parameters for the Delete method.
type DeleteRequest struct {
	Post model.Post `json:"post"`
}

// DeleteResponse collects the response parameters for the Delete method.
type DeleteResponse struct {
	Response string `json:"response"`
	Err      error  `json:"err"`
}

// MakeDeleteEndpoint returns an endpoint that invokes Delete on the service.
func MakeDeleteEndpoint(s service.PostsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		response, err := s.Delete(ctx, req.Post)
		return DeleteResponse{
			Err:      err,
			Response: response,
		}, nil
	}
}

// Failed implements Failer.
func (r DeleteResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Store implements Service. Primarily useful in a client.
func (e Endpoints) Store(ctx context.Context, post model.Post) (response string, err error) {
	request := StoreRequest{Post: post}
	response0, err := e.StoreEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response0.(StoreResponse).Response, response0.(StoreResponse).Err
}

// Update implements Service. Primarily useful in a client.
func (e Endpoints) Update(ctx context.Context, post model.Post) (response string, err error) {
	request := UpdateRequest{Post: post}
	response0, err := e.UpdateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response0.(UpdateResponse).Response, response0.(UpdateResponse).Err
}

// List implements Service. Primarily useful in a client.
func (e Endpoints) List(ctx context.Context, post model.Post) (response []*pb.Post, err error) {
	request := ListRequest{Post: post}
	response0, err := e.ListEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response0.(ListResponse).Response, response0.(ListResponse).Err
}

// Delete implements Service. Primarily useful in a client.
func (e Endpoints) Delete(ctx context.Context, post model.Post) (response string, err error) {
	request := DeleteRequest{Post: post}
	response0, err := e.DeleteEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response0.(DeleteResponse).Response, response0.(DeleteResponse).Err
}
