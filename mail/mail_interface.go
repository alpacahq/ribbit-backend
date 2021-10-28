package mail

import "github.com/alpacahq/ribbit-backend/model"

// Service is the interface to access our Mail
type Service interface {
	Send(subject string, toName string, toEmail string, content string, HTMLContent string) error
	SendWithDefaults(subject, toEmail, content string, HTMLContent string) error
	SendVerificationEmail(toEmail string, v *model.Verification) error
	SendForgotVerificationEmail(toEmail string, v *model.Verification) error
}
