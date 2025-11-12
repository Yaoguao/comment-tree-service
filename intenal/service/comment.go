package service

import (
	"comment-tree-service/intenal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

var InvalidID = errors.New("invalid ID")

type CommentsStorage interface {
	Create(ctx context.Context, c *models.Comment) error
	GetThread(ctx context.Context, parentID bson.ObjectID, limit, offset int, sort string) ([]models.Comment, error)
	DeleteThread(ctx context.Context, id bson.ObjectID) error
	Search(ctx context.Context, query string, limit, offset int) ([]models.Comment, error)
}

type CommentsService struct {
	Storage CommentsStorage
}

func NewCommentsService(storage CommentsStorage) *CommentsService {
	return &CommentsService{Storage: storage}
}

func (s *CommentsService) SaveComment(c *models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return s.Storage.Create(ctx, c)
}

func (s *CommentsService) GetThread(parentIDHex string, limit, offset int, sort string) ([]models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	parentID, err := bson.ObjectIDFromHex(parentIDHex)
	if err != nil {
		return nil, InvalidID
	}

	return s.Storage.GetThread(ctx, parentID, limit, offset, sort)
}

func (s *CommentsService) DeleteThread(IDHex string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	id, err := bson.ObjectIDFromHex(IDHex)
	if err != nil {
		return InvalidID
	}

	return s.Storage.DeleteThread(ctx, id)
}
func (s *CommentsService) Search(query string, limit, offset int) ([]models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return s.Storage.Search(ctx, query, limit, offset)
}
