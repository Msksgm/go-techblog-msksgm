package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/msksgm/go-techblog-msksgm/model"
)

func articleResponse(article *model.Article) M {
	if article == nil {
		return nil
	}
	return M{
		"title":     article.Title,
		"body":      article.Body,
		"slug":      article.Slug,
		"createdAt": article.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt": article.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (s *Server) createArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title string `json:"title" validate:"required"`
			Body  string `json:"body" validate:"required"`
			Slug  string `json:"slug" validate:"required"`
		} `json:"article"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Article); err != nil {
			validationError(w, err)
			return
		}

		article := model.Article{
			Title: input.Article.Title,
			Body:  input.Article.Body,
			Slug:  input.Article.Slug,
		}

		user := userFromContext(r.Context())
		article.Author = user
		article.AuthorID = user.ID

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		if err := s.articleService.CreateArticle(r.Context(), &article); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, M{"article": article})
	}
}

func (s *Server) listArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := model.ArticleFilter{}

		if v := query.Get("author"); v != "" {
			filter.AuthorUsername = &v
		}

		articles, err := s.articleService.Articles(r.Context(), filter)
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"articles": articles})
	}
}

func (s *Server) getArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filter := model.ArticleFilter{}

		if slug, exists := vars["slug"]; exists {
			filter.Slug = &slug
		}

		articles, err := s.articleService.Articles(r.Context(), filter)
		if err != nil {
			serverError(w, err)
			return
		}

		var article *model.Article

		if len(articles) > 0 {
			article = articles[0]
		}

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) updateArticle() http.HandlerFunc {
	type Input struct {
		Article struct {
			Title *string `json:"title,omitempty"`
			Body  *string `json:"body,omitempty"`
		} `json:"article,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		slug := mux.Vars(r)["slug"]

		article, err := s.articleService.ArticleBySlug(r.Context(), slug)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrNotFound):
				err := ErrorM{"article": []string{"requested article not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())

		if user.ID != article.AuthorID {
			err := ErrorM{"article": []string{"forbidden request"}}
			errorResponse(w, http.StatusForbidden, err)
			return
		}

		patch := model.ArticlePatch{
			Title: input.Article.Title,
			Body:  input.Article.Body,
		}

		if err := s.articleService.UpdateArticle(r.Context(), article, patch); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"article": article})
	}
}

func (s *Server) deleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]

		article, err := s.articleService.ArticleBySlug(r.Context(), slug)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrNotFound):
				err := ErrorM{"article": []string{"requested article not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())

		if user.ID != article.AuthorID {
			err := ErrorM{"article": []string{"forbidden request"}}
			errorResponse(w, http.StatusForbidden, err)
			return
		}

		if err := s.articleService.DeleteArticle(r.Context(), article.ID); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusNoContent, nil)
	}
}
