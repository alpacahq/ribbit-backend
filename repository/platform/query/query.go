package query

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
)

// List prepares data for list queries
func List(u *model.AuthUser) (*model.ListQuery, error) {
	switch true {
	case int(u.Role) <= 2: // user is SuperAdmin or Admin
		return nil, nil
	default:
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
}
