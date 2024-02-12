package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct {
	Message string
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. HealthzResponse を変数に代入
	res := &model.HealthzResponse{Message: "OK"} // 1. Message の Field に OK という文字を入れる

	// 2. JSONにシリアライズして書き込み
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(res) // 3. Encodeメソッドを呼び出し
	if err != nil {
		// 4. シリアライズ失敗時のエラーログ出力
		log.Println("Failed to encode HealthzResponse:", err)
	}
}
