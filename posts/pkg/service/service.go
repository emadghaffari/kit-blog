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

	"github.com/emadghaffari/kit-blog/posts/pkg/grpc/pb"
	model "github.com/emadghaffari/kit-blog/posts/pkg/model"
)

// PostsService describes the service.
type PostsService interface {
	// Add your methods here
	Store(ctx context.Context, post model.Post) (response string, err error)
	Update(ctx context.Context, post model.Post) (response string, err error)
	List(ctx context.Context, post model.Post) (response []*pb.Post, err error)
	Delete(ctx context.Context, post model.Post) (response string, err error)
}

type basicPostsService struct {
	db *mongo.Collection
}

func (b *basicPostsService) Store(ctx context.Context, post model.Post) (response string, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("store")
	defer span.Finish()

	values := bson.M{
		"title":       post.Title,
		"slug":        post.Slug,
		"description": post.Description,
		"body":        post.Body,
		"header":      post.Header,
		"createdAt":   post.CreatedAT,
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

	return oid.Hex(), nil
}
func (b *basicPostsService) Update(ctx context.Context, post model.Post) (response string, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("update")
	defer span.Finish()

	oid, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return "FAILD", err
	}

	filter := bson.M{"_id": oid}
	values := bson.M{
		"title":       post.Title,
		"slug":        post.Slug,
		"description": post.Description,
		"body":        post.Body,
		"header":      post.Header,
		"createdAt":   post.CreatedAT,
	}
	_, err = b.db.ReplaceOne(context.Background(), filter, values)
	if err != nil {
		return "FAILD", nil
	}

	return oid.Hex(), err
}
func (b *basicPostsService) List(ctx context.Context, post model.Post) (response []*pb.Post, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("list")
	defer span.Finish()

	items := []*pb.Post{}
	cur, err := b.db.Find(context.Background(), bson.M{})
	if err != nil {
		return items, nil
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &model.Post{}
		err := cur.Decode(data)
		if err != nil {
			return items, nil
		}
		items = append(items,
			&pb.Post{
				Title:       data.Title,
				Body:        data.Body,
				Slug:        data.Slug,
				Description: data.Description,
				CreatedAt:   &pb.Post_Time{Time: *data.CreatedAT},
				Header:      data.Header,
			},
		)
	}

	return items, err
}
func (b *basicPostsService) Delete(ctx context.Context, post model.Post) (response string, err error) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("delete")
	defer span.Finish()

	oid, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return "FAILD", err
	}

	filter := bson.M{"_id": oid}
	_, err = b.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return "FAILD", err
	}

	return oid.Hex(), err
}

// NewBasicPostsService returns a naive, stateless implementation of PostsService.
func NewBasicPostsService() PostsService {
	col, err := initMongoDB()
	if err != nil {
		return new(basicPostsService)
	}
	return &basicPostsService{
		db: col,
	}
}

// New returns a PostsService with all of the expected middleware wired in.
func New(middleware []Middleware) PostsService {
	var svc PostsService = NewBasicPostsService()
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

	posts := client.Database("kit-posts").Collection("posts")

	return posts, nil

}
