package mockgopg

import (
	"strings"
	"sync"

	"github.com/alpacahq/ribbit-backend/manager"
)

// SQLMock handles query mocks
type SQLMock struct {
	lock          *sync.RWMutex
	currentQuery  string // tracking queries
	currentParams []interface{}
	queries       map[string]buildQuery
	currentInsert string // tracking inserts
	inserts       map[string]buildInsert
}

// ExpectInsert is a builder method that accepts a model as interface and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectInsert(models ...interface{}) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	var inserts []string
	for _, v := range models {
		inserts = append(inserts, strings.ToLower(manager.GetType(v)))
	}
	currentInsert := strings.Join(inserts, ",")

	sqlMock.currentInsert = currentInsert
	return sqlMock
}

// ExpectExec is a builder method that accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectExec(query string) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = strings.TrimSpace(query)
	return sqlMock
}

// ExpectQuery accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectQuery(query string) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = strings.TrimSpace(query)
	return sqlMock
}

// ExpectQueryOne accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) ExpectQueryOne(query string) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = strings.TrimSpace(query)
	return sqlMock
}

// WithArgs is a builder method that accepts a query in string and returns an SQLMock pointer
func (sqlMock *SQLMock) WithArgs(params ...interface{}) *SQLMock {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentParams = make([]interface{}, 0)
	for _, p := range params {
		sqlMock.currentParams = append(sqlMock.currentParams, p)
	}

	return sqlMock
}

// Returns accepts expected result and error, and completes the build of our sqlMock object
func (sqlMock *SQLMock) Returns(result *OrmResult, err error) {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	q := buildQuery{
		query:  sqlMock.currentQuery,
		params: sqlMock.currentParams,
		result: result,
		err:    err,
	}
	sqlMock.queries[sqlMock.currentQuery] = q
	sqlMock.currentQuery = ""
	sqlMock.currentParams = nil

	i := buildInsert{
		insert: sqlMock.currentInsert,
		err:    err,
	}
	sqlMock.inserts[sqlMock.currentInsert] = i
	sqlMock.currentInsert = ""
}

// FlushAll resets our sqlMock object
func (sqlMock *SQLMock) FlushAll() {
	sqlMock.lock.Lock()
	defer sqlMock.lock.Unlock()

	sqlMock.currentQuery = ""
	sqlMock.currentParams = nil
	sqlMock.queries = make(map[string]buildQuery)

	sqlMock.currentInsert = ""
	sqlMock.inserts = make(map[string]buildInsert)
}
