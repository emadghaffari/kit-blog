// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package grpc

import (
	endpoint "github.com/emadghaffari/kit-blog/users/pkg/endpoint"
	pb "github.com/emadghaffari/kit-blog/users/pkg/grpc/pb"
	grpc "github.com/go-kit/kit/transport/grpc"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	get      grpc.Handler
	login    grpc.Handler
	register grpc.Handler
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.UsersServer {
	return &grpcServer{
		get:      makeGetHandler(endpoints, options["Get"]),
		login:    makeLoginHandler(endpoints, options["Login"]),
		register: makeRegisterHandler(endpoints, options["Register"]),
	}
}
