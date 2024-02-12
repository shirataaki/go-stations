package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	// HealthzHandlerのエンドポイントを登録
	healthzHandler := handler.NewHealthzHandler() // HealthzHandlerのインスタンスを作成
	mux.Handle("/healthz", healthzHandler)        // /healthz のエンドポイントに healthzHandler を割り当て

	return mux
}
