package mockgopg

import "github.com/go-pg/pg/v9/orm"

// OrmResult struct to implements orm.Result
type OrmResult struct {
	rowsAffected int
	rowsReturned int
	model        interface{}
}

// Model implements an orm.Model
func (o *OrmResult) Model() orm.Model {
	if o.model == nil {
		return nil
	}

	model, err := orm.NewModel(o.model)
	if err != nil {
		return nil
	}

	return model
}

// RowsAffected returns the number of rows affected in the data table
func (o *OrmResult) RowsAffected() int {
	return o.rowsAffected
}

// RowsReturned returns the number of rows
func (o *OrmResult) RowsReturned() int {
	return o.rowsReturned
}

// NewResult implements orm.Result in go-pg package
func NewResult(rowAffected, rowReturned int, model interface{}) *OrmResult {
	return &OrmResult{
		rowsAffected: rowAffected,
		rowsReturned: rowReturned,
		model:        model,
	}
}
