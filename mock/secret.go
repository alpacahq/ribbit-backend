package mock

// Password mock
type Password struct {
	HashPasswordFn        func(string) string
	HashMatchesPasswordFn func(hash, password string) bool
	HashRandomPasswordFn  func() (string, error)
}

// HashPassword mock
func (p *Password) HashPassword(password string) string {
	return p.HashPasswordFn(password)
}

// HashMatchesPassword mock
func (p *Password) HashMatchesPassword(hash, password string) bool {
	return p.HashMatchesPasswordFn(hash, password)
}

// HashRandomPassword mock
func (p *Password) HashRandomPassword() (string, error) {
	return p.HashRandomPasswordFn()
}
