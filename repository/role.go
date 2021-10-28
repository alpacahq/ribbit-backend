package repository

import (
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

// NewRoleRepo returns a Role Repo instance
func NewRoleRepo(db *pg.DB, log *zap.Logger) *RoleRepo {
	return &RoleRepo{db, log}
}

// RoleRepo represents the client for the role table
type RoleRepo struct {
	db  *pg.DB
	log *zap.Logger
}

// CreateRoles creates role objects in our database
func (r *RoleRepo) CreateRoles() error {
	role := new(model.Role)
	sql := `INSERT INTO roles (id, access_level, name) VALUES (?, ?, ?) ON CONFLICT DO NOTHING`
	r.db.Query(role, sql, 1, model.SuperAdminRole, "superadmin")
	r.db.Query(role, sql, 2, model.AdminRole, "admin")
	r.db.Query(role, sql, 3, model.UserRole, "user")
	return nil
}
