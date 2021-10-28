package mockgopg

import "github.com/go-pg/pg/v9/orm"

// Formatter implements orm.Formatter
type Formatter struct {
}

// FormatQuery formats our query and params to byte
func (f *Formatter) FormatQuery(b []byte, query string, params ...interface{}) []byte {
	formatter := new(orm.Formatter)
	got := formatter.FormatQuery(b, query, params...)
	return got
}
