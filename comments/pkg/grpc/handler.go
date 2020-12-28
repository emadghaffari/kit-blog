package grpc

import (
	"context"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"

	endpoint "github.com/emadghaffari/kit-blog/comments/pkg/endpoint"
	pb "github.com/emadghaffari/kit-blog/comments/pkg/grpc/pb"
	"github.com/emadghaffari/kit-blog/comments/pkg/service"
)

// makeStoreHandler creates the handler logic
func makeStoreHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.StoreEndpoint, decodeStoreRequest, encodeStoreResponse, options...)
}

// decodeStoreResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Store request.
func decodeStoreRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StoreRequest)
	return endpoint.StoreRequest{Cm: service.Comment{UserID: req.UserID, PostID: req.PostID, Title: req.Title, Body: req.Body}}, nil
}

// encodeStoreResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeStoreResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.StoreResponse)
	if resp.Err != nil {
		return &pb.StoreReply{Id: "", Status: pb.StoreReply_Fail.String()}, resp.Err
	}
	return &pb.StoreReply{Id: resp.Id, Status: pb.StoreReply_Success.String()}, nil
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
	return endpoint.UpdateRequest{Cm: service.Comment{UserID: req.UserID, PostID: req.PostID, Title: req.Title, Body: req.Body, ID: req.Id}}, nil

}

// encodeUpdateResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.UpdateResponse)
	if resp.Err != nil {
		return &pb.UpdateReply{Id: "", Status: pb.UpdateReply_Fail.String()}, resp.Err
	}
	return &pb.UpdateReply{Id: resp.Id, Status: pb.UpdateReply_Success.String()}, nil
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
	return endpoint.ListRequest{PostID: req.PostID}, nil

}

// encodeListResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeListResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.ListResponse)
	if resp.Err != nil {
		return &pb.ListReply{Comments: []*pb.Comment{}, Status: pb.ListReply_Fail.String()}, resp.Err
	}
	return &pb.ListReply{Comments: resp.CMS, Status: pb.ListReply_Success.String()}, nil
}
func (g *grpcServer) List(ctx context1.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	_, rep, err := g.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListReply), nil
}
