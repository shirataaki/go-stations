package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	// HealthzHandlerのエンドポイントを登録
	healthzHandler := handler.NewHealthzHandler() // HealthzHandlerのインスタンスを作成
	mux.Handle("/healthz", healthzHandler)        // /healthz のエンドポイントに healthzHandler を割り当て

	todoService := service.NewTODOService(todoDB)      // TODOServiceのインスタンスを作成
	todoHandler := handler.NewTODOHandler(todoService) // TODOHandlerのインスタンスを作成
	mux.Handle("/todos", todoHandler)

	// 必ずpanicを発生させるHandler
	/*mux.Handle("/do-panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic")
	}))*/

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional panic")
	})
	mux.Handle("/do-panic", middleware.Recovery(panicHandler))

	return mux
}
