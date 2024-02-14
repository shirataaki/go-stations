package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

//var ErrNotFound = errors.New("not found")

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	if id <= 0 {
		return nil, &model.ErrNotFound{}
	}

	if subject == "" {
		// Return a generic error for simplicity, in real scenario, this should be validated before reaching the DB layer.
		return nil, errors.New("subject is required") // Adjust error handling as per your application logic
	}

	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
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
		return nil, &model.ErrNotFound{}
	}

	// Retrieve the updated TODO
	var todo model.TODO
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &model.ErrNotFound{}
		}
		return nil, fmt.Errorf("failed to retrieve updated todo: %w", err)
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
