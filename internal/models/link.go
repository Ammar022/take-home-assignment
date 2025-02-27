package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Link represents a link in the bio
type Link struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" binding:"required"`
	URL       string             `bson:"url" json:"url" binding:"required,url"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
	Clicks    int                `bson:"clicks" json:"clicks"`
	UserID    string             `bson:"userId" json:"userId"`
}

// LinkCreateDTO is used for creating a new link
type LinkCreateDTO struct {
	Title     string    `json:"title" binding:"required"`
	URL       string    `json:"url" binding:"required,url"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    string    `json:"userId"`
}

// LinkUpdateDTO is used for updating an existing link
type LinkUpdateDTO struct {
	Title     string    `json:"title"`
	URL       string    `json:"url" binding:"omitempty,url"`
	ExpiresAt time.Time `json:"expiresAt"`
}

