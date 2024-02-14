package model

import "fmt"

// ErrNotFound は指定されたリソースが見つからない場合のエラーを表す。
type ErrNotFound struct {
	Resource string // 見つからなかったリソースの種類
	ID       int64  // 見つからなかったリソースのID
}

// ErrNotFound 構造体の Error メソッドを定義
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}
