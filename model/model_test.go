package model_test

import (
	"testing"

	"github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/model"
)

func TestBeforeInsert(t *testing.T) {
	base := &model.Base{}
	base.BeforeInsert(nil)
	if base.CreatedAt.IsZero() {
		t.Errorf("CreatedAt was not changed")
	}
	if base.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt was not changed")
	}
}

func TestBeforeUpdate(t *testing.T) {
	base := &model.Base{
		CreatedAt: mock.TestTime(2000),
	}
	base.BeforeUpdate(nil)
	if base.UpdatedAt == mock.TestTime(2001) {
		t.Errorf("UpdatedAt was not changed")
	}

}

func TestDelete(t *testing.T) {
	baseModel := &model.Base{
		CreatedAt: mock.TestTime(2000),
		UpdatedAt: mock.TestTime(2001),
	}
	baseModel.Delete()
	if baseModel.DeletedAt.IsZero() {
		t.Errorf("DeletedAt not changed")
	}

}
