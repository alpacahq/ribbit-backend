package manager

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/gertd/go-pluralize"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// NewManager returns a new manager
func NewManager(accountRepo *repository.AccountRepo, roleRepo *repository.RoleRepo, db *pg.DB) *Manager {
	return &Manager{accountRepo, roleRepo, db}
}

// Manager holds a group of methods for writing tests
type Manager struct {
	accountRepo *repository.AccountRepo
	roleRepo    *repository.RoleRepo
	db          *pg.DB
}

// CreateSchema creates tables declared as models (struct)
func (m *Manager) CreateSchema(models ...interface{}) {
	for _, model := range models {
		opt := &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}
		err := m.db.CreateTable(model, opt)
		if err != nil {
			log.Fatal(err)
		}
		p := pluralize.NewClient()
		modelName := GetType(model)
		tableName := p.Plural(strings.ToLower(modelName))
		fmt.Printf("Created model %s as table %s\n", modelName, tableName)
	}
}

// CreateRoles is a thin wrapper for roleRepo.CreateRoles(), which populates our roles table
func (m *Manager) CreateRoles() {
	err := m.roleRepo.CreateRoles()
	if err != nil {
		log.Fatal(err)
	}
}

// CreateSuperAdmin is used to create a user object with superadmin role
func (m *Manager) CreateSuperAdmin(email, password string) (*model.User, error) {
	u := &model.User{
		Email:    email,
		Password: secret.New().HashPassword(password),
		Active:   true,
		Verified: true,
		RoleID:   int(model.SuperAdminRole),
	}
	return m.accountRepo.Create(u)
}

// GetType is a useful utility function to help us inspect the name of a model (struct) which is expressed as an interface{}
func GetType(myvar interface{}) string {
	valueOf := reflect.ValueOf(myvar)
	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	}
	return valueOf.Type().Name()
}

// GetModels retrieve models
func GetModels() []interface{} {
	return model.Models
}
