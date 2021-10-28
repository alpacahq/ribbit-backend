package model_test

import (
	"testing"

	"github.com/alpacahq/ribbit-backend/model"
)

func TestUpdateLastLogin(t *testing.T) {
	user := &model.User{
		FirstName: "TestGuy",
	}
	user.UpdateLastLogin()
	if user.LastLogin.IsZero() {
		t.Errorf("Last login time was not changed")
	}
}

func TestUpdateUpdatedAt(t *testing.T) {
	user := &model.User{
		FirstName: "TestGal",
	}
	user.Update()
	if user.UpdatedAt.IsZero() {
		t.Errorf("updated_at is not changed")
	}
}

func TestUpdateDeletedAt(t *testing.T) {
	user := &model.User{
		FirstName: "TestGod",
	}
	user.Delete()
	if user.DeletedAt.IsZero() {
		t.Errorf("deleted_at is not changed")
	}
}
