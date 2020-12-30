package endpoint

import (
	"context"

	endpoint "github.com/go-kit/kit/endpoint"

	service "github.com/emadghaffari/kit-blog/users/pkg/service"
)

// LoginRequest collects the request parameters for the Login method.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse collects the response parameters for the Login method.
type LoginResponse struct {
	S0 string `json:"s0"`
	E1 error  `json:"e1"`
}

// MakeLoginEndpoint returns an endpoint that invokes Login on the service.
func MakeLoginEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		s0, e1 := s.Login(ctx, req.Username, req.Password)
		return LoginResponse{
			E1: e1,
			S0: s0,
		}, nil
	}
}

// Failed implements Failer.
func (r LoginResponse) Failed() error {
	return r.E1
}

// RegisterRequest collects the request parameters for the Register method.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// RegisterResponse collects the response parameters for the Register method.
type RegisterResponse struct {
	S0 string `json:"s0"`
	E1 error  `json:"e1"`
}

// MakeRegisterEndpoint returns an endpoint that invokes Register on the service.
func MakeRegisterEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		s0, e1 := s.Register(ctx, req.Username, req.Password, req.Email, req.Phone)
		return RegisterResponse{
			E1: e1,
			S0: s0,
		}, nil
	}
}

// Failed implements Failer.
func (r RegisterResponse) Failed() error {
	return r.E1
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Login implements Service. Primarily useful in a client.
func (e Endpoints) Login(ctx context.Context, username string, password string) (s0 string, e1 error) {
	request := LoginRequest{
		Password: password,
		Username: username,
	}
	response, err := e.LoginEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(LoginResponse).S0, response.(LoginResponse).E1
}

// Register implements Service. Primarily useful in a client.
func (e Endpoints) Register(ctx context.Context, username string, password string, email string, phone string) (s0 string, e1 error) {
	request := RegisterRequest{
		Email:    email,
		Password: password,
		Phone:    phone,
		Username: username,
	}
	response, err := e.RegisterEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RegisterResponse).S0, response.(RegisterResponse).E1
}

// GetRequest collects the request parameters for the Get method.
type GetRequest struct {
	Id string `json:"id"`
}

// GetResponse collects the response parameters for the Get method.
type GetResponse struct {
	S0 string `json:"username"`
	S1 string `json:"email"`
	S2 string `json:"phone"`
	E1 error  `json:"error"`
}

// MakeGetEndpoint returns an endpoint that invokes Get on the service.
func MakeGetEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		s0, s1, s2, e1 := s.Get(ctx, req.Id)
		return GetResponse{
			E1: e1,
			S0: s0,
			S1: s1,
			S2: s2,
		}, nil
	}
}

// Failed implements Failer.
func (r GetResponse) Failed() error {
	return r.E1
}

// Get implements Service. Primarily useful in a client.
func (e Endpoints) Get(ctx context.Context, id string) (s0 string, s1 string, s2 string, e1 error) {
	request := GetRequest{Id: id}
	response, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).S0, response.(GetResponse).S1, response.(GetResponse).S2, response.(GetResponse).E1
}
