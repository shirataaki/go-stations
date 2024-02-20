package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// TODOをDBに保存
	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	// 保存したTODOのIDを取得
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todo id: %w", err)
	}

	// 保存したTODOを読み取り
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todo: %w", err)
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	// サイズが0の場合は空のスライスを返す
	if size <= 0 {
		return []*model.TODO{}, nil
	}

	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var rows *sql.Rows
	var err error

	if prevID > 0 {
		log.Printf("Executing クエリ with prevID: %v, size: %v\n", prevID, size)
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		log.Printf("Executing query with size: %v\n", size)
		rows, err = s.db.QueryContext(ctx, read, size)
	}

	if err != nil {
		log.Printf("Query execution error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var todos []*model.TODO

	for rows.Next() {
		var todo model.TODO
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Printf("エラー scanning row: %v\n", err)
			return nil, err
		}
		log.Printf("スキャン TODO: ID=%d, Subject=%s\n", todo.ID, todo.Subject)
		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return nil, err
	}

	if todos == nil { // nilの返却を避けるために空のスライスにする。
		log.Printf("todosはnilです。")
		todos = []*model.TODO{}
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	/*if id <= 0 {
		return nil, &model.ErrNotFound{}
	}*/

	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?` // idを含めるように変更した
	)

	// Execute the update query
	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check affected rows: %w", err)
	}

	if affected == 0 {
		return nil, &model.ErrNotFound{} // 更新された行がない場合は、ErrNotFoundエラーを返す
	}

	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// IDに対応するTODOが見つからない場合は、ErrNotFoundエラーを具体的な情報と共に返す
			return nil, &model.ErrNotFound{Resource: "TODO", ID: id}
		}
		return nil, fmt.Errorf("failed to retrieve updated todo: %w", err)
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	// idsが空のスライスの場合は何もせずに終了
	if len(ids) == 0 {
		return nil
	}

	// 削除クエリのWHERE句で使用するプレースホルダーを生成
	placeholder := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf(`DELETE FROM todos WHERE id IN (%s)`, placeholder)

	// int64のスライスをinterface{}のスライスに変換
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// 削除クエリの実行
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete todos: %w", err)
	}

	// 削除された行がない場合はErrNotFoundを返す
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if affected == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
