package grpc

import (
	"context"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"

	"github.com/emadghaffari/kit-blog/posts/model"
	endpoint "github.com/emadghaffari/kit-blog/posts/pkg/endpoint"
	pb "github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
)

// makeStoreHandler creates the handler logic
func makeStoreHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.StoreEndpoint, decodeStoreRequest, encodeStoreResponse, options...)
}

// decodeStoreResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Store request.
func decodeStoreRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StoreRequest)
	return endpoint.StoreRequest{
		Post: model.Post{
			Token:       &req.Post.Token,
			Title:       req.Post.Title,
			Body:        req.Post.Body,
			Slug:        req.Post.Slug,
			Description: req.Post.Description,
			Header:      req.Post.Header,
		},
	}, nil
}

// encodeStoreResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeStoreResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.StoreResponse)
	if resp.Err != nil {
		return &pb.StoreReply{
			Response: "",
			Status:   pb.StoreReply_Fail,
		}, nil
	}
	return &pb.StoreReply{
		Response: resp.Response,
		Status:   pb.StoreReply_Success,
	}, nil
}
func (g *grpcServer) Store(ctx context1.Context, req *pb.StoreRequest) (*pb.StoreReply, error) {
	_, rep, err := g.store.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.StoreReply), nil
}

// makeUpdateHandler creates the handler logic
func makeUpdateHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdateEndpoint, decodeUpdateRequest, encodeUpdateResponse, options...)
}

// decodeUpdateResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Update request.
func decodeUpdateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UpdateRequest)

	return endpoint.UpdateRequest{
		Post: model.Post{
			Token:       &req.Post.Token,
			Title:       req.Post.Title,
			Body:        req.Post.Body,
			Slug:        req.Post.Slug,
			Description: req.Post.Description,
			Header:      req.Post.Header,
		},
	}, nil
}

// encodeUpdateResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.UpdateResponse)
	if resp.Err != nil {
		return &pb.UpdateReply{
			Response: "",
			Status:   pb.UpdateReply_Fail,
		}, nil
	}
	return &pb.UpdateReply{
		Response: resp.Response,
		Status:   pb.UpdateReply_Success,
	}, nil
}
func (g *grpcServer) Update(ctx context1.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	_, rep, err := g.update.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateReply), nil
}

// makeListHandler creates the handler logic
func makeListHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.ListEndpoint, decodeListRequest, encodeListResponse, options...)
}

// decodeListResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain List request.
func decodeListRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ListRequest)

	return endpoint.ListRequest{
		Post: model.Post{
			Token:       &req.Post.Token,
			Title:       req.Post.Title,
			Body:        req.Post.Body,
			Slug:        req.Post.Slug,
			Description: req.Post.Description,
			Header:      req.Post.Header,
		},
	}, nil
}

// encodeListResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeListResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.ListResponse)
	if resp.Err != nil {
		return &pb.ListReply{
			Post: []*pb.Post{},
		}, nil
	}
	return &pb.ListReply{
		Post: resp.Response,
	}, nil
}
func (g *grpcServer) List(ctx context1.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	_, rep, err := g.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListReply), nil
}

// makeDeleteHandler creates the handler logic
func makeDeleteHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.DeleteEndpoint, decodeDeleteRequest, encodeDeleteResponse, options...)
}

// decodeDeleteResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Delete request.
func decodeDeleteRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.DeleteRequest)

	return endpoint.DeleteRequest{
		Post: model.Post{
			Token:       &req.Post.Token,
			Title:       req.Post.Title,
			Body:        req.Post.Body,
			Slug:        req.Post.Slug,
			Description: req.Post.Description,
			Header:      req.Post.Header,
		},
	}, nil
}

// encodeDeleteResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeDeleteResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.DeleteResponse)
	if resp.Err != nil {
		return &pb.DeleteReply{
			Response: resp.Response,
			Status:   pb.DeleteReply_Fail,
		}, nil
	}
	return &pb.DeleteReply{
		Response: resp.Response,
		Status:   pb.DeleteReply_Success,
	}, nil
}
func (g *grpcServer) Delete(ctx context1.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	_, rep, err := g.delete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DeleteReply), nil
}
