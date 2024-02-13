package model

import "time"

type (
	// A TODO expresses ...
	TODO struct {
		ID          int64     `json:"id"`
		Subject     string    `json:"subject"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// 利用者から受け取る値の定義
	CreateTODORequest struct {
		Subject     string `json:"subject" binding:"required"`
		Description string `json:"description"` // 必須ではない
	}
	// 利用者に返す値（この場合は構造体）の定義
	CreateTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)
