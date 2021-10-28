package model

import (
	"context"
	"time"
)

// Models hold registered models in-memory
var Models []interface{}

// Base contains common fields for all tables
type Base struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Pagination holds pagination's data
type Pagination struct {
	Limit  int
	Offset int
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    int
}

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}

// Delete sets deleted_at time to current_time
func (b *Base) Delete() {
	t := time.Now()
	b.DeletedAt = &t
}

// Register is used for registering models
func Register(m interface{}) {
	Models = append(Models, m)
}
