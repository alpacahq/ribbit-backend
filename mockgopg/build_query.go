package mockgopg

type buildQuery struct {
	funcName string
	query    string
	params   []interface{}
	result   *OrmResult
	err      error
}
