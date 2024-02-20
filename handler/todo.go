package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// CreateTODORequest に JSON Decode
		var req model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// subject が空文字列の場合を判定
		if req.Subject == "" {
			http.Error(w, "Bad Request: subject is required", http.StatusBadRequest)
			return
		}

		// CreateTODO メソッドを呼び出し
		todo, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// CreateTODOResponse に代入し、JSON Encode を行い HTTP Response を返す
		if todo != nil {
			resp := model.CreateTODOResponse{TODO: *todo} // ポインタから値を取り出す
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			// todoがnilの場合の適切なエラーハンドリング
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPut:
		// UpdateTODORequest に JSON Decode
		var req model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// id が 0 の場合や subject が空文字列の場合を判定
		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request: ID and subject are required", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil {
			var notFound *model.ErrNotFound
			if errors.As(err, &notFound) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		// 更新成功時のレスポンスを返す
		resp := model.UpdateTODOResponse{TODO: *todo}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	case http.MethodGet:
		// クエリパラメータからprev_idとsizeを取得し、整数に変換
		prevIDStr := r.URL.Query().Get("prev_id")
		sizeStr := r.URL.Query().Get("size")
		var prevID int64 = 0 // デフォルト値として0を設定
		var size int64 = 10  // デフォルト値として10を設定（例）

		if prevIDStr != "" {
			var err error
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid prev_id parameter", http.StatusBadRequest)
				return
			}
		}

		if sizeStr != "" {
			var err error
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid size parameter", http.StatusBadRequest)
				return
			}
		}

		// ReadTODOメソッドを呼び出してTODOリストを取得
		todos, err := h.svc.ReadTODO(r.Context(), prevID, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// レスポンスにContent-Typeを設定
		w.Header().Set("Content-Type", "application/json")
		// 取得したTODOリストをエンコードしてレスポンスとして返す
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"todos": todos}); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		// DeleteTODORequestにJSON Decode
		var req model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// idのリストが空の場合を判定
		if len(req.IDs) == 0 {
			http.Error(w, "Bad Request: ids are required", http.StatusBadRequest)
			return
		}

		// DeleteTODOメソッドを呼び出し
		err = h.svc.DeleteTODO(r.Context(), req.IDs)
		if err != nil {
			// ErrNotFoundが返却された場合は404 NotFoundとしてHTTP Responseを返す
			var notFound *model.ErrNotFound
			if errors.As(err, &notFound) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				// その他のエラーの場合は500 Internal Server Errorを返す
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// 削除できた場合は、DeleteTODOResponseを作成し、JSON Encodeを行いHTTP Responseを返す
		resp := model.DeleteTODOResponse{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	default:
		// PUTとPOST以外のメソッドに対してはMethod Not Allowedを返す
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
