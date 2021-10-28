package mail

import (
	"os"
	"strings"

	"github.com/alpacahq/ribbit-backend/config"
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/sendgrid/sendgrid-go"
	s "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// NewMail generates new Mail variable
func NewMail(mc *config.MailConfig, sc *config.SiteConfig) *Mail {
	return &Mail{
		ExternalURL: sc.ExternalURL,
		FromName:    "Ribbit App",
		FromEmail:   "ribbitapp@ribbitapp.co",
	}
}

// Mail provides a mail service implementation
type Mail struct {
	ExternalURL string
	FromName    string
	FromEmail   string
}

// Send email with sendgrid
func (m *Mail) Send(subject string, toName string, toEmail string, content string, HTMLContent string) error {
	from := s.NewEmail(m.FromName, m.FromEmail)
	to := s.NewEmail(toName, toEmail)
	message := s.NewSingleEmail(from, subject, to, content, HTMLContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

// SendWithDefaults assumes some defaults for sending out email with sendgrid
func (m *Mail) SendWithDefaults(subject, toEmail, content string, HTMLContent string) error {
	err := m.Send(subject, toEmail, toEmail, content, HTMLContent)
	if err != nil {
		return err
	}
	return nil
}

// SendVerificationEmail assumes defaults and generates a verification email
func (m *Mail) SendVerificationEmail(toEmail string, v *model.Verification) error {
	// url := m.ExternalURL + "/verification/" + v.Token
	content := "Here is your otp: " + v.Token
	HTMLContent := `<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office"><head><meta charset="UTF-8"/><meta http-equiv="X-UA-Compatible" content="IE=edge"/><meta http-equiv="Content-Type" content="text/html charset=UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1.0"/><meta name="x-apple-disable-message-reformatting"/><title>Alpaca Email</title><link rel="preconnect" href="https://fonts.gstatic.com"/><link rel="preconnect" href="https://fonts.gstatic.com"/><link href="https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap" rel="stylesheet"/><linkhref="https://fonts.googleapis.com/css2?family=Roboto+Mono:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;1,100;1,200;1,300;1,400;1,500;1,600;1,700&family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap"rel="stylesheet"/><linkrel="stylesheet"href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.3/css/all.min.css"integrity="sha512-iBBXm8fW90+nuLcSKlbmrPcLa0OT92xO1BIsZ+ywDWZCvqsWgccV3gFoRBv0z+8dLJgyAHIhR35VZc2oM/gI1w=="crossorigin="anonymous"referrerpolicy="no-referrer"/><!--[if mso]><noscript><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml></noscript><![endif]--><style>table,td,div,h1,p{font-family: "Roboto", sans-serif;}</style></head><body style="margin: 0; padding: 0;"><table role="presentation" style="width: 100%; border-collapse: collapse; border: 0; border-spacing: 0; background: #ffffff;"><tr><td align="center" style="padding: 0;"><table role="presentation" style="max-width: 602px; border-collapse: collapse; border-spacing: 0; text-align: left;"><tr><td align="center" style="padding: 5% 0 5% 0;"><img src="http://35.193.43.181:8080/file/assets/img/header_logo.png" alt="company-logo" width="220" style="height: auto; display: block;"/></td></tr><tr><td style="padding: 36px 30px 42px 30px;"><table role="presentation" style="width: 100%; border-collapse: collapse; border: 0; border-spacing: 0;"><tr><td align="center" style="padding: 0 0 20% 0; color: #153643; border-bottom: 1px solid #cddde7;"><img src="http://35.193.43.181:8080/file/assets/img/verify_email.png" alt="verify_email" width="120" height="120" style="display: block; margin: 20px 0 40px 0;"/><h1 style="font-family: 'Roboto', sans-serif; font-size: 28px; color: rgba(53, 55, 80, 1); margin: 5% 0 5% 0; padding: 0 15px 0 15px;">Verify your email</h1><p style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; color: rgba(53, 64, 80, 0.49); margin: 10% 0 10% 0; padding: 0 15px 0 15px;">Enter this code to verify your email.</p><table role="presentation" style="width: 100%; border-collapse: collapse; border: 0; border-spacing: 0; text-align: center;"><tr><td style="width: auto; padding: 0; font-size: 0; line-height: 0;">&nbsp;</td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_1</p></td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_2</p></td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_3</p></td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_4</p></td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_5</p></td><td align="center" style="max-width: 30px; padding: 0 4px; vertical-align: center; color: #3d4a52; border-bottom: 1px solid #e8e8e8; display: inline-block; margin: 0 10px;"><p align="center" style="font-family: 'Roboto', sans-serif; font-size: 21px; line-height: 24px; margin: 15px 0;">TOKEN_6</p></td><td style="width: auto; padding: 0; font-size: 0; line-height: 0;">&nbsp;</td></tr></table></td></tr></table></td></tr><tr><td style="padding: 5%; background: #ffffff;"><table role="presentation" style="width: 100%; border-collapse: collapse; border: 0; border-spacing: 0;"><tr><td align="center" style="padding: 0; width: 100%;"><p style="font-family: 'Roboto', sans-serif; font-size: 14px; line-height: 16px; color: #74787a; margin: 8% 0 5% 0; padding: 0 35px 0 35px;">760 Market Street, Floor 10 San Francisco, CA, 94102</p></td><td style="padding: 0; width: 50%;" align="right"></td></tr></table></td></tr></table></td></tr></table></body></html>`
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_1", string(v.Token[0]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_2", string(v.Token[1]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_3", string(v.Token[2]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_4", string(v.Token[3]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_5", string(v.Token[4]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_6", string(v.Token[5]), -1)
	err := m.SendWithDefaults("Verification Email", toEmail, content, HTMLContent)
	if err != nil {
		return err
	}
	return nil
}

// SendForgotVerificationEmail assumes defaults and generates a verification email
func (m *Mail) SendForgotVerificationEmail(toEmail string, v *model.Verification) error {
	// url := m.ExternalURL + "/verification/" + v.Token
	content := "Here is your otp: " + v.Token
	HTMLContent := `<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office"><head> <meta charset="UTF-8"> <meta http-equiv="X-UA-Compatible" content="IE=edge"> <meta http-equiv="Content-Type" content="text/html charset=UTF-8"/> <meta name="viewport" content="width=device-width, initial-scale=1.0"> <meta name="x-apple-disable-message-reformatting"> <title>Alpaca Email</title> <link rel="preconnect" href="https://fonts.gstatic.com"> <link rel="preconnect" href="https://fonts.gstatic.com"> <link href="https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap" rel="stylesheet"> <link href="https://fonts.googleapis.com/css2?family=Roboto+Mono:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;1,100;1,200;1,300;1,400;1,500;1,600;1,700&family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap" rel="stylesheet"> <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.3/css/all.min.css" integrity="sha512-iBBXm8fW90+nuLcSKlbmrPcLa0OT92xO1BIsZ+ywDWZCvqsWgccV3gFoRBv0z+8dLJgyAHIhR35VZc2oM/gI1w==" crossorigin="anonymous" referrerpolicy="no-referrer"/><!--[if mso]><noscript><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml></noscript><![endif]--> <style>table, td, div, h1, p{font-family: 'Roboto', sans-serif;}</style></head><body style="margin:0;padding:0;"><table role="presentation" style="width:100%;border-collapse:collapse;border:0;border-spacing:0;background:#ffffff;"><tr><td align="center" style="padding:0;"><table role="presentation" style="max-width:602px; border-collapse:collapse; border-spacing:0; text-align:left;"> <tr><td align="center" style=" padding:5% 0 5% 0;"><img src="http://35.193.43.181:8080/file/assets/img/header_logo.png" alt="company-logo" width="220" style="height:auto;display:block;"/></td></tr><tr><td style="padding:36px 30px 42px 30px;"><table role="presentation" style="width:100%; border-collapse:collapse; border:0; border-spacing:0;"><tr><td align="center" style="padding:0 0 20% 0;color:#153643; border-bottom: 1px solid #CDDDE7;"> <img src="http://35.193.43.181:8080/file/assets/img/forgot_password.png" alt="verify_email" width="120" height="120" style="display:block; margin: 20px 0 40px 0"/><h1 style="font-family: 'Roboto', sans-serif; font-size:28px; color:rgba(53, 55, 80, 1); margin:5% 0 5% 0; padding: 0 15px 0 15px; "> Forgot Password? </h1><p style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; color: rgba(53, 64, 80, 0.49); margin:10% 0 10% 0; padding: 0 15px 0 15px;"> Enter this code to reset your password. </p><table role="presentation" style="width:100%; border-collapse:collapse; border:0; border-spacing:0; text-align: center;"> <tr> <td style="width:auto; padding:0; font-size:0; line-height:0;">&nbsp;</td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_1 </p></td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_2 </p></td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_3 </p></td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_4 </p></td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_5 </p></td><td align="center" style="max-width:30px; padding: 0 4px; vertical-align:center; color: #3D4A52; border-bottom: 1px solid #E8E8E8; display: inline-block; margin: 0 10px;"> <p align="center" style="font-family: 'Roboto', sans-serif; font-size:21px; line-height:24px; margin: 15px 0; "> TOKEN_6 </p></td><td style="width:auto; padding:0; font-size:0; line-height:0;">&nbsp;</td></tr></table></td></tr></table></td></tr><tr><td style="padding:5%; background:#ffffff"><table role="presentation" style="width:100%; border-collapse:collapse; border:0; border-spacing:0;"><tr><td align="center" style="padding:0; width:100%;" > <p style="font-family: 'Roboto', sans-serif; font-size:14px; line-height:16px; color:#74787A; margin:8% 0 5% 0; padding: 0 35px 0 35px;">760 Market Street, Floor 10 San Francisco, CA, 94102</p></td><td style="padding:0;width:50%;" align="right"> </td></tr></table></td></tr></table></td></tr></table></body></html>`
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_1", string(v.Token[0]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_2", string(v.Token[1]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_3", string(v.Token[2]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_4", string(v.Token[3]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_5", string(v.Token[4]), -1)
	HTMLContent = strings.Replace(HTMLContent, "TOKEN_6", string(v.Token[5]), -1)

	err := m.SendWithDefaults("Verification Email", toEmail, content, HTMLContent)
	if err != nil {
		return err
	}
	return nil
}
