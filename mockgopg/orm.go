package mockgopg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/alpacahq/ribbit-backend/manager"

	"github.com/go-pg/pg/v9/orm"
)

type goPgDB struct {
	sqlMock *SQLMock
}

// NewGoPGDBTest returns method that already implements orm.DB and mock instance to mocking arguments and results.
func NewGoPGDBTest() (conn orm.DB, mock *SQLMock, err error) {
	sqlMock := &SQLMock{
		lock:          new(sync.RWMutex),
		currentQuery:  "",
		currentParams: nil,
		queries:       make(map[string]buildQuery),
		currentInsert: "",
		inserts:       make(map[string]buildInsert),
	}

	goPG := &goPgDB{
		sqlMock: sqlMock,
	}

	return goPG, sqlMock, nil
}

// not yet implemented
func (p *goPgDB) Model(model ...interface{}) *orm.Query {
	return nil
}

func (p *goPgDB) ModelContext(c context.Context, model ...interface{}) *orm.Query {
	return nil
}

func (p *goPgDB) Select(model interface{}) error {
	return nil
}

func (p *goPgDB) Insert(model ...interface{}) error {
	// return nil
	return p.doInsert(context.Background(), model...)
}

func (p *goPgDB) Update(model interface{}) error {
	return nil
}

func (p *goPgDB) Delete(model interface{}) error {
	return nil
}

func (p *goPgDB) ForceDelete(model interface{}) error {
	return nil
}

func (p *goPgDB) Exec(query interface{}, params ...interface{}) (orm.Result, error) {
	sqlQuery := fmt.Sprintf("%v", query)
	return p.doQuery(context.Background(), nil, sqlQuery, params...)
}

func (p *goPgDB) ExecContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	sqlQuery := fmt.Sprintf("%v", query)
	return p.doQuery(c, nil, sqlQuery, params...)
}

func (p *goPgDB) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (p *goPgDB) ExecOneContext(c context.Context, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (p *goPgDB) Query(model, query interface{}, params ...interface{}) (orm.Result, error) {
	sqlQuery := fmt.Sprintf("%v", query)
	return p.doQuery(context.Background(), model, sqlQuery, params...)
}

func (p *goPgDB) QueryContext(c context.Context, model, query interface{}, params ...interface{}) (orm.Result, error) {
	sqlQuery := fmt.Sprintf("%v", query)
	return p.doQuery(c, model, sqlQuery, params...)
}

func (p *goPgDB) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	sqlQuery := fmt.Sprintf("%v", query)
	return p.doQuery(context.Background(), model, sqlQuery, params...)
}

func (p *goPgDB) QueryOneContext(c context.Context, model, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (p *goPgDB) CopyFrom(r io.Reader, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (p *goPgDB) CopyTo(w io.Writer, query interface{}, params ...interface{}) (orm.Result, error) {
	return nil, nil
}

func (p *goPgDB) Context() context.Context {
	return context.Background()
}

func (p *goPgDB) Formatter() orm.QueryFormatter {
	f := new(Formatter)
	return f
}

func (p *goPgDB) doInsert(ctx context.Context, models ...interface{}) error {
	// update p.insertMock
	for k, v := range p.sqlMock.inserts {

		// not handling value at the moment

		onTheListInsertStr := k

		var inserts []string
		for _, v := range models {
			inserts = append(inserts, strings.ToLower(manager.GetType(v)))
		}
		wantedInsertStr := strings.Join(inserts, ",")

		if onTheListInsertStr == wantedInsertStr {
			return v.err
		}
	}

	return nil
}

func (p *goPgDB) doQuery(ctx context.Context, dst interface{}, query string, params ...interface{}) (orm.Result, error) {
	// replace duplicate space
	space := regexp.MustCompile(`\s+`)

	for k, v := range p.sqlMock.queries {
		onTheList := p.Formatter().FormatQuery(nil, k, v.params...)
		onTheListQueryStr := strings.TrimSpace(space.ReplaceAllString(string(onTheList), " "))

		wantedQuery := p.Formatter().FormatQuery(nil, query, params...)
		wantedQueryStr := strings.TrimSpace(space.ReplaceAllString(string(wantedQuery), " "))

		if onTheListQueryStr == wantedQueryStr {
			var (
				data []byte
				err  error
			)

			if dst == nil {
				return v.result, v.err
			}

			data, err = json.Marshal(v.result.model)
			if err != nil {
				return v.result, err
			}

			err = json.Unmarshal(data, dst)
			if err != nil {
				return v.result, err
			}

			return v.result, v.err
		}
	}

	return nil, fmt.Errorf("no mock expectation result")
}
