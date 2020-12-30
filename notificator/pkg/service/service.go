package service

import (
	"context"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/emadghaffari/api-teacher/utils/random"
)

// NotificatorService describes the service.
type NotificatorService interface {
	// Add your methods here
	Send(ctx context.Context, to, body string) (id string, err error)
}

type basicNotificatorService struct {
	db *mongo.Collection
}

func (b *basicNotificatorService) Send(ctx context.Context, to string, body string) (id string, err error) {
	values := bson.M{
		"to":    to,
		"body":  body,
		"notif": random.Rand(10000, 999999),
	}
	res, err := b.db.InsertOne(context.Background(), values)

	if err != nil {
		log.Printf("Error in insert data to mongodb: %v", err)
		return "Fail", err
	}

	if _, ok := res.InsertedID.(primitive.ObjectID); !ok {
		log.Printf("Error in get oid from res")
		return "Fail", err
	}

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			span := tracer.StartSpan("notification", opentracing.ChildOf(pctx))
			defer span.Finish()

			span.SetTag("sended notification to", to)
		}
	}

	return "SUCCESS", err
}

// NewBasicNotificatorService returns a naive, stateless implementation of NotificatorService.
func NewBasicNotificatorService() NotificatorService {

	col, err := initMongoDB()
	if err != nil {
		return new(basicNotificatorService)
	}
	return &basicNotificatorService{
		db: col,
	}
}

// New returns a NotificatorService with all of the expected middleware wired in.
func New(middleware []Middleware) NotificatorService {
	var svc NotificatorService = NewBasicNotificatorService()
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

	users := client.Database("kit-notification").Collection("notification")

	return users, nil

}
