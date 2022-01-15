package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func Test_healthcheck(t *testing.T) {
	// テストサーバーの作成
	srv := testServer()
	// テストのリクエストを定義とレスポンスを初期化
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	// テストのリクエストを実行
	srv.router.ServeHTTP(w, req)
	gotResp := M{}
	extractResponseBody(w.Body, &gotResp)
	// レスポンスの期待値を宣言
	expectResp := M{
		"message": "health",
		"status":  "available",
	}

	// ステータスコードを比較
	if code := w.Code; code != http.StatusOK {
		t.Errorf("expected status code of 200, but got %d", code)
	}

	// レスポンスの期待値を実数値を比較
	if !reflect.DeepEqual(expectResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectResp, gotResp)
	}
}

func testServer() *Server {
	srv := &Server{
		router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

func extractResponseBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		panic(err)
	}
}
