package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/msksgm/go-techblog-msksgm/mock"
	"github.com/msksgm/go-techblog-msksgm/model"
)

func Test_createArticle(t *testing.T) {
	articleStore := &mock.ArticleService{}
	userStore := &mock.UserService{}
	srv := testServer()
	srv.articleService = articleStore
	srv.userService = userStore

	token, err := generateUserToken(
		&model.User{
			ID:       1,
			Username: "username",
		},
	)
	if err != nil {
		panic(err)
	}

	input := `{
		"article": {
			"title": "title",
			"body": "body",
			"slug": "slug"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", strings.NewReader(input))
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))
	w := httptest.NewRecorder()

	currentUser := &model.User{
		Username: "username",
		Token:    token,
	}
	userStore.GetCurrentUserFn = func() *model.User {
		return currentUser
	}

	var article model.Article
	articleStore.CreateArticleFn = func(a *model.Article) error {
		article = *a
		return nil
	}
	srv.router.ServeHTTP(w, req)
	expectedResp := articleResponse(&article)

	gotResp := M{}
	extractResponseArticleBody(w.Body, &gotResp)

	if code := w.Code; code != http.StatusCreated {
		t.Errorf("expected status code of 201, but got %d", code)
	}

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

func extractResponseArticleBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm["article"])
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		panic(err)
	}
}
