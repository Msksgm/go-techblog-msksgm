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

type ArticleFilter struct {
	ID             *uint
	Title          *string
	AuthorID       *uint
	AuthorUsername *string
	Slug           *string

	Limit  int
	Offset int
}

type ArticlePatch struct {
	Title *string
	Body  *string
	Slug  *string
}

type ArticleService interface {
	CreateArticle(context.Context, *Article) error
	ArticleBySlug(context.Context, string) (*Article, error)
	Articles(context.Context, ArticleFilter) ([]*Article, error)
	UpdateArticle(context.Context, *Article, ArticlePatch) error
	DeleteArticle(context.Context, uint) error
}
