package mock

import "github.com/alpacahq/ribbit-backend/model"

// Mail mock
type Mail struct {
	ExternalURL             string
	SendFn                  func(string, string, string, string) error
	SendWithDefaultsFn      func(string, string, string) error
	SendVerificationEmailFn func(string, *model.Verification) error
}

// Send mock
func (m *Mail) Send(subject, toName, toEmail, content string) error {
	return m.SendFn(subject, toName, toEmail, content)
}

// SendWithDefaults mock
func (m *Mail) SendWithDefaults(subject, toEmail, content string) error {
	return m.SendWithDefaultsFn(subject, toEmail, content)
}

// SendVerificationEmail mock
func (m *Mail) SendVerificationEmail(toEmail string, v *model.Verification) error {
	return m.SendVerificationEmailFn(toEmail, v)
}
