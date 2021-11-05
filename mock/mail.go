package mock

import "github.com/alpacahq/ribbit-backend/model"

// Mail mock
type Mail struct {
	ExternalURL                   string
	SendFn                        func(string, string, string, string, string) error
	SendWithDefaultsFn            func(string, string, string, string) error
	SendVerificationEmailFn       func(string, *model.Verification) error
	SendForgotVerificationEmailFn func(string, *model.Verification) error
}

// Send mock
func (m *Mail) Send(subject, toName, toEmail, content, html string) error {
	return m.SendFn(subject, toName, toEmail, content, html)
}

// SendWithDefaults mock
func (m *Mail) SendWithDefaults(subject, toEmail, content, html string) error {
	return m.SendWithDefaultsFn(subject, toEmail, content, html)
}

// SendVerificationEmail mock
func (m *Mail) SendVerificationEmail(toEmail string, v *model.Verification) error {
	return m.SendVerificationEmailFn(toEmail, v)
}

func (m *Mail) SendForgotVerificationEmail(toEmail string, v *model.Verification) error {
	return m.SendForgotVerificationEmailFn(toEmail, v)
}
