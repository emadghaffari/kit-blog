package endpoint

import (
	"context"
	service "github.com/emadghaffari/kit-blog/notificator/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// SendRequest collects the request parameters for the Send method.
type SendRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

// SendResponse collects the response parameters for the Send method.
type SendResponse struct {
	Id  string `json:"id"`
	Err error  `json:"err"`
}

// MakeSendEndpoint returns an endpoint that invokes Send on the service.
func MakeSendEndpoint(s service.NotificatorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendRequest)
		id, err := s.Send(ctx, req.To, req.Body)
		return SendResponse{
			Err: err,
			Id:  id,
		}, nil
	}
}

// Failed implements Failer.
func (r SendResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Send implements Service. Primarily useful in a client.
func (e Endpoints) Send(ctx context.Context, to string, body string) (id string, err error) {
	request := SendRequest{
		Body: body,
		To:   to,
	}
	response, err := e.SendEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SendResponse).Id, response.(SendResponse).Err
}
