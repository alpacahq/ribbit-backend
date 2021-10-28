package mock

// Mobile mock
type Mobile struct {
	GenerateSMSTokenFn func(string, string) error
	CheckCodeFn        func(string, string, string) error
}

// GenerateSMSToken mock
func (m *Mobile) GenerateSMSToken(countryCode, mobile string) error {
	return m.GenerateSMSTokenFn(countryCode, mobile)
}

// CheckCode mock
func (m *Mobile) CheckCode(countryCode, mobile, code string) error {
	return m.CheckCodeFn(countryCode, mobile, code)
}
