package mock

import (
	"context"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type ArticleService struct {
	CreateArticleFn func(*model.Article) error
	ArticleBySlugFn func(*model.Article) (*model.Article, error)
	ArticlesFn      func() ([]*model.Article, error)
}

func (m *ArticleService) CreateArticle(_ context.Context, article *model.Article) error {
	return m.CreateArticleFn(article)
}

func (m *ArticleService) Articles(_ context.Context, af model.ArticleFilter) ([]*model.Article, error) {
	return m.ArticlesFn()
}

func (m *ArticleService) ArticleBySlug(_ context.Context, slug string) (*model.Article, error) {
	return nil, nil
}

func (m *ArticleService) DeleteArticle(_ context.Context, id uint) error {
	return nil
}

func (m *ArticleService) UpdateArticle(_ context.Context, article *model.Article, patch model.ArticlePatch) error {
	return nil
}
