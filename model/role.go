package model

func init() {
	Register(&Role{})
}

// AccessRole represents access role type
type AccessRole int8

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = iota + 1

	// AdminRole has admin specific permissions
	AdminRole

	// UserRole is a standard user
	UserRole
)

// Role model
type Role struct {
	ID          int        `json:"id"`
	AccessLevel AccessRole `json:"access_level"`
	Name        string     `json:"name"`
}

// RoleRepo represents the database interface
type RoleRepo interface {
	CreateRoles() error
}
