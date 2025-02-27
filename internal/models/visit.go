package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Visit represents a click on a link
type Visit struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LinkID    primitive.ObjectID `bson:"linkId" json:"linkId"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	UserAgent string             `bson:"userAgent" json:"userAgent"`
	IP        string             `bson:"ip" json:"ip"`
	Referrer  string             `bson:"referrer" json:"referrer"`
}