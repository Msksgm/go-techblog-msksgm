package model

import (
	"context"
	"time"
)

type Article struct {
	ID        uint      `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Slug      string    `json:"slug"`
	AuthorID  uint      `json:"-" db:"author_id"`
	Author    *User     `json:"-"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ArticleService interface {
	CreateArticle(context.Context, *Article) error
}
