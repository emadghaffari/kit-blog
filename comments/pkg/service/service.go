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

	"github.com/emadghaffari/kit-blog/comments/config"
	"github.com/emadghaffari/kit-blog/comments/pkg/grpc/pb"
	us "github.com/emadghaffari/kit-blog/users/pkg/grpc/pb"
)

// Comment struct
type Comment struct {
	ID     string `json:"id,omitempty" bson:"id"`
	UserID string `json:"user_id,omitempty" bson:"user_id"`
	PostID string `json:"post_id,omitempty" bson:"post_id"`
	Title  string `json:"title,omitempty" bson:"title"`
	Body   string `json:"body,omitempty" bson:"body"`
}

// CommentsService describes the service.
type CommentsService interface {
	// Add your methods here
	Store(ctx context.Context, cm Comment) (id string, err error)
	Update(ctx context.Context, cm Comment) (id string, err error)
	List(ctx context.Context, postID string) (cms []*pb.Comment, err error)
}

type basicCommentsService struct {
	user us.UsersClient
	db   *mongo.Collection
}

func (b *basicCommentsService) Store(ctx context.Context, cm Comment) (id string, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("store")

	values := bson.M{
		"post_id": cm.PostID,
		"user_id": cm.UserID,
		"title":   cm.Title,
		"body":    cm.Body,
	}
	res, err := b.db.InsertOne(context.Background(), values)

	if err != nil {
		log.Printf("Error in insert data to mongodb: %v", err)
		return "FAILD", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error in get oid from res")
		return "FAILD", err
	}

	span.Finish()
	return oid.Hex(), nil
}

func (b *basicCommentsService) Update(ctx context.Context, cm Comment) (id string, err error) {
	oid, err := primitive.ObjectIDFromHex(cm.ID)
	if err != nil {
		return "FAILD", err
	}

	data := Comment{}
	filter := bson.M{"_id": oid}
	res := b.db.FindOne(context.Background(), filter)
	if err := res.Decode(&data); err != nil {
		return "FAILD", err
	}
	data.Body = cm.Body
	data.Title = cm.Title

	_, err = b.db.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return "FAILD", nil
	}

	return oid.Hex(), err
}

func (b *basicCommentsService) List(ctx context.Context, postID string) (cms []*pb.Comment, err error) {

	items := []*pb.Comment{}
	cur, err := b.db.Find(context.Background(), bson.M{"post_id": postID})
	if err != nil {
		return items, nil
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &Comment{}
		err := cur.Decode(data)
		if err != nil {
			return items, nil
		}
		res, err := b.user.Get(context.Background(), &us.GetRequest{Id: data.UserID})
		if err != nil {
			return items, nil
		}
		items = append(items,
			&pb.Comment{PostID: data.PostID,
				Title:     data.Title,
				Body:      data.Body,
				Username:  &pb.Comment_Name{Name: res.Username},
				Useremail: &pb.Comment_Email{Email: res.Email},
			},
		)
	}

	return items, err
}

// NewBasicCommentsService returns a naive, stateless implementation of CommentsService.
func NewBasicCommentsService() CommentsService {
	conn, err := initUsers()
	if err != nil {
		return new(basicCommentsService)
	}

	col, err := initMongoDB()
	if err != nil {
		return new(basicCommentsService)
	}

	return &basicCommentsService{
		user: us.NewUsersClient(conn),
		db:   col,
	}
}

// New returns a CommentsService with all of the expected middleware wired in.
func New(middleware []Middleware) CommentsService {
	var svc CommentsService = NewBasicCommentsService()
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

	users := client.Database("kit-comments").Collection("comments")

	return users, nil

}

func initUsers() (*grpc.ClientConn, error) {
	// Users from etcd database
	tracer := opentracing.GlobalTracer()
	conn, err := grpc.Dial(config.Confs.Users.Host,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())))
	if err != nil {
		log.Printf("unable to connect to users service, %s", err.Error())
		return nil, err
	}
	return conn, nil
}
