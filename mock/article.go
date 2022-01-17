package mock

import (
	"context"

	"github.com/msksgm/go-techblog-msksgm/model"
)

type ArticleService struct {
	CreateArticleFn func(*model.Article) error
}

func (m *ArticleService) CreateArticle(_ context.Context, article *model.Article) error {
	return m.CreateArticleFn(article)
}

func (m *ArticleService) Articles(_ context.Context, af model.ArticleFilter) ([]*model.Article, error) {
	return nil, nil
}
