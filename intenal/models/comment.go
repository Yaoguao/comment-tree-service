package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Comment struct {
	ID        bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	ParentID  *bson.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Path      string         `bson:"path,omitempty" json:"path"`
	Content   string         `bson:"content" json:"content"`
	Author    string         `bson:"author" json:"author"`
	CreatedAt time.Time      `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at,omitempty" json:"updated_at"`
}
