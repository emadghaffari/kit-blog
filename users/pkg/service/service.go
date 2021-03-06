package service

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	cryptoutils "github.com/emadghaffari/api-teacher/utils/cryptoUtils"
	"github.com/emadghaffari/kit-blog/notificator/pkg/grpc/pb"
	"github.com/emadghaffari/kit-blog/users/config"
	"github.com/emadghaffari/kit-blog/users/pkg/model"
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
	db                *mongo.Collection
}

func (b *basicUsersService) Login(ctx context.Context, username string, password string) (s0 string, e1 error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("login")
	defer span.Finish()

	data := model.User{}
	values := bson.M{"username": username, "password": cryptoutils.GetMD5(password)}
	res := b.db.FindOne(context.Background(), values)
	if res.Err() != nil {
		return "", res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return "", err
	}

	// send notification service
	ct := opentracing.ContextWithSpan(context.Background(), span)
	if _, err := b.notificatorClient.Send(ct, &pb.SendRequest{To: data.Phone, Body: "Hi " + username}); err != nil {
		log.Printf("failed to send notif: %v", err)
		return "", err
	}

	jwt, err := model.Conf.Generate(data)
	if err != nil {
		log.Printf("Error in create jwt: %v", err)
		return "", err
	}

	return jwt.AccessToken, err
}
func (b *basicUsersService) Register(ctx context.Context, username string, password string, email string, phone string) (s0 string, e1 error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("register")
	defer span.Finish()

	values := bson.M{
		"username": username,
		"password": cryptoutils.GetMD5(password),
		"email":    email,
		"phone":    phone,
	}
	res, err := b.db.InsertOne(context.Background(), values)

	if err != nil {
		log.Printf("Error in insert data to mongodb: %v", err)
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error in get oid from res")
		return "", err
	}

	// send notification service
	ct := opentracing.ContextWithSpan(context.Background(), span)
	if _, err = b.notificatorClient.Send(ct, &pb.SendRequest{To: phone, Body: "Hi " + username}); err != nil {
		log.Printf("failed to send notif: %v", err)
		return "", err
	}

	jwt, err := model.Conf.Generate(model.User{ID: oid.Hex(), Username: username, Email: email, Phone: phone})
	if err != nil {
		log.Printf("Error in create jwt: %v", err)
		return "", err
	}

	return jwt.AccessToken, err
}

func (b *basicUsersService) Get(ctx context.Context, id string) (username, email, phone string, err error) {

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			span := tracer.StartSpan("get_user", opentracing.ChildOf(pctx))
			defer span.Finish()
		}
	}

	user := model.User{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", "", "", err
	}
	values := bson.M{"_id": oid}
	res := b.db.FindOne(context.Background(), values)
	if err := res.Decode(&user); err != nil {
		return "", "", "", err
	}

	return user.Username, user.Email, user.Phone, nil
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {
	conn, err := initNotificator()
	if err != nil {
		return new(basicUsersService)
	}

	col, err := initMongoDB()
	if err != nil {
		return new(basicUsersService)
	}

	return &basicUsersService{
		notificatorClient: pb.NewNotificatorClient(conn),
		db:                col,
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

func initMongoDB() (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@mongo:27017"))
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	users := client.Database("kit-users").Collection("users")

	return users, nil

}

func initNotificator() (*grpc.ClientConn, error) {
	// notificator from etcd database
	tracer := opentracing.GlobalTracer()
	conn, err := grpc.Dial(config.Confs.Notifs.Host,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		log.Printf("unable to connect to notificator service, %s", err.Error())
		return nil, err
	}
	return conn, nil
}
