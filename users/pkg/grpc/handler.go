package grpc

import (
	"context"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"

	endpoint "github.com/emadghaffari/kit-blog/users/pkg/endpoint"
	pb "github.com/emadghaffari/kit-blog/users/pkg/grpc/pb"
)

func makeLoginHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.LoginEndpoint, decodeLoginRequest, encodeLoginResponse, options...)
}

func decodeLoginRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return endpoint.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

func encodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.LoginResponse)
	if resp.E1 != nil {
		return &pb.LoginReply{Token: "", Status: pb.LoginReply_Fail}, resp.E1
	}
	return &pb.LoginReply{Token: resp.S0, Status: pb.LoginReply_Success}, nil
}
func (g *grpcServer) Login(ctx context1.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	_, rep, err := g.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginReply), nil
}

func makeRegisterHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.RegisterEndpoint, decodeRegisterRequest, encodeRegisterResponse, options...)
}

func decodeRegisterRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.RegisterRequest)
	return endpoint.RegisterRequest{Username: req.Username, Password: req.Password, Email: req.Email, Phone: req.Phone}, nil
}

func encodeRegisterResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.RegisterResponse)
	if resp.E1 != nil {
		return &pb.RegisterReply{Token: "", Status: pb.RegisterReply_Fail}, resp.E1
	}
	return &pb.RegisterReply{Token: resp.S0, Status: pb.RegisterReply_Success}, nil
}
func (g *grpcServer) Register(ctx context1.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	_, rep, err := g.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RegisterReply), nil
}

func makeGetHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.GetEndpoint, decodeGetRequest, encodeGetResponse, options...)
}

func decodeGetRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetRequest)
	return endpoint.GetRequest{Id: req.Id}, nil
}

func encodeGetResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.GetResponse)
	if resp.E1 != nil {
		return &pb.GetReply{Username: "", Email: "", Phone: "", Status: pb.GetReply_Fail}, resp.E1
	}
	return &pb.GetReply{Username: resp.S0, Email: resp.S1, Phone: resp.S2, Status: pb.GetReply_Success}, nil
}
func (g *grpcServer) Get(ctx context1.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetReply), nil
}
