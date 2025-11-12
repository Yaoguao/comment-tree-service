package mongodb

import (
	"comment-tree-service/intenal/config"
	"comment-tree-service/intenal/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type MongoStorage struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewMongoStorage(cfg *config.Config) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.Storage.MongoDB.DSN)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database(cfg.Storage.MongoDB.DBName)
	col := db.Collection("comments")

	// индексы
	_, _ = col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "content", Value: "text"}},
		},
		{
			Keys: bson.D{{Key: "path", Value: 1}},
		},
	})

	return &MongoStorage{
		client: client,
		col:    col,
	}, nil
}

func (m *MongoStorage) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoStorage) Create(ctx context.Context, c *models.Comment) error {
	c.ID = bson.NewObjectID()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	if c.ParentID != nil {
		var parent models.Comment
		err := m.col.FindOne(ctx, bson.M{"_id": c.ParentID}).Decode(&parent)
		if err != nil {
			return fmt.Errorf("parent not found: %w", err)
		}
		c.Path = parent.Path + "/" + c.ID.Hex()
	} else {
		c.Path = c.ID.Hex()
	}

	_, err := m.col.InsertOne(ctx, c)
	return err
}

func (m *MongoStorage) GetThread(ctx context.Context, parentID bson.ObjectID, limit, offset int, sort string) ([]models.Comment, error) {
	var parent models.Comment
	err := m.col.FindOne(ctx, bson.M{"_id": parentID}).Decode(&parent)
	if err != nil {
		return nil, fmt.Errorf("parent not found: %w", err)
	}

	filter := bson.M{"path": bson.M{"$regex": "^" + parent.Path}}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	if sort == "asc" {
		opts.SetSort(bson.D{{Key: "created_at", Value: 1}})
	} else {
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *MongoStorage) DeleteThread(ctx context.Context, id bson.ObjectID) error {
	var parent models.Comment
	if err := m.col.FindOne(ctx, bson.M{"_id": id}).Decode(&parent); err != nil {
		return err
	}
	_, err := m.col.DeleteMany(ctx, bson.M{"path": bson.M{"$regex": "^" + parent.Path}})
	return err
}

func (m *MongoStorage) Search(ctx context.Context, query string, limit, offset int) ([]models.Comment, error) {
	filter := bson.M{}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	if query != "" {
		filter = bson.M{"$text": bson.M{"$search": query}}
		opts.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	} else {
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // например сортировка по дате
	}

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}
