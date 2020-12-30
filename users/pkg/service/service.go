package service

import (
	"context"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	etcd "go.etcd.io/etcd/client/v2"
	"google.golang.org/grpc"

	"github.com/emadghaffari/kit-blog/notificator/pkg/grpc/pb"
)

// UsersService describes the service.
type UsersService interface {
	// Add your methods here
	Get(ctx context.Context, id string) (username, email, phone string, err error)
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password, email, phone string) (string, error)
}

type basicUsersService struct {
	notificatorClient pb.NotificatorClient
}

func (b *basicUsersService) Login(ctx context.Context, username string, password string) (s0 string, e1 error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("login")
	defer span.Finish()

	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}

	ct := opentracing.ContextWithSpan(context.Background(), span)
	_, err = b.notificatorClient.Send(ct, &pb.SendRequest{To: "09355960597", Body: "Hi Welcome."})
	if err != nil {
		log.Printf("failed to send notif: %v", err)
	}

	return u.String(), err
}
func (b *basicUsersService) Register(ctx context.Context, username string, password string, email string, phone string) (s0 string, e1 error) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}
	return u.String(), err
}

func (b *basicUsersService) Get(ctx context.Context, id string) (username, email, phone string, err error) {
	// TODO implement the business logic of Get
	return "Emad", "emadghaffariii@gmail.com", "09355980597", nil
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {
	var (
		prefix = "/blog/notificator/notificator"
	)
	cfg := etcd.Config{
		Endpoints: []string{"http://etcd:2379"},
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := etcd.New(cfg)
	if err != nil {
		log.Printf("unable to connect to etcd: %s", err.Error())
		return new(basicUsersService)
	}
	kapi := etcd.NewKeysAPI(c)
	en, err := kapi.Get(context.Background(), prefix, &etcd.GetOptions{})
	if err != nil {
		log.Printf("unable to get entries from etcd: %s", err)
		return new(basicUsersService)
	}

	// notificator from etcd database
	log.Printf("---------------%v-----------------", en.Node.Value)
	tracer := opentracing.GlobalTracer()
	conn, err := grpc.Dial("localhost:8082",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		log.Printf("unable to connect to notificator service, %s", err.Error())
		return new(basicUsersService)
	}
	return &basicUsersService{
		notificatorClient: pb.NewNotificatorClient(conn),
	}
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	var svc UsersService = NewBasicUsersService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
