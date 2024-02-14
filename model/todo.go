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
	ReadTODORequest struct {
		PrevID int64 `json:"prev_id"` // 前回取得した最後のTODOのID
		Size   int64 `json:"size"`
	}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	UpdateTODORequest struct {
		ID          int64  `json:"id" binding:"required"`      // 必須
		Subject     string `json:"subject" binding:"required"` // 必須
		Description string `json:"description"`                // 必須ではない
	}

	UpdateTODOResponse struct {
		TODO TODO `json:"todo"` // 変更されたTODO
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)
