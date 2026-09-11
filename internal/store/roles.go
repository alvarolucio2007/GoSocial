package store

import (
	"context"
	"database/sql"
)

type RolesRepository interface {
	GetByName(ctx context.Context, name string) (*Role, error)
}

type RoleLevel int16

const (
	LevelUser RoleLevel = iota + 1
	LevelMod
	LevelAdmin
)

type Role struct {
	ID          int64     `json:"role_id"`
	Name        string    `json:"role_name"`
	Level       RoleLevel `json:"role_level"`
	Description string    `json:"role_description"`
}
type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id,name,description,level FROM roles WHERE name=$1`
	role := &Role{}
	err := s.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}
	return role, nil
}
