package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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
