package grpc

import (
	"context"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"

	endpoint "github.com/emadghaffari/kit-blog/notificator/pkg/endpoint"
	pb "github.com/emadghaffari/kit-blog/notificator/pkg/grpc/pb"
)

// makeSendHandler creates the handler logic
func makeSendHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SendEndpoint, decodeSendRequest, encodeSendResponse, options...)
}

// decodeSendResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Send request.
func decodeSendRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendRequest)

	return endpoint.SendRequest{To: req.To, Body: req.Body}, nil
}

// encodeSendResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSendResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.SendResponse)
	if resp.Err != nil {
		return &pb.SendReply{Id: "", Status: pb.SendReply_Fail}, resp.Err
	}
	return &pb.SendReply{Id: resp.Id, Status: pb.SendReply_Success}, nil
}
func (g *grpcServer) Send(ctx context1.Context, req *pb.SendRequest) (*pb.SendReply, error) {
	_, rep, err := g.send.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendReply), nil
}
