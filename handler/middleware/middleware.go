package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

// 独自の型を定義
type contextKey string

// コンテキストに格納するためのキーとして使用する独自の型の変数を定義
var contextKeyOSName = contextKey("OSName")

func UserAgentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		osName := ua.OS

		// 独自の型のキーを使用してOS名をコンテキストに格納
		ctx := context.WithValue(r.Context(), contextKeyOSName, osName)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
