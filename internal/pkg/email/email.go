package email

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendMailV3 for send mail from sendgrid
func SendMailV3(from *mail.Email, p *mail.Personalization, templateID string) error {
	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.SetTemplateID(templateID)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)

	return err
}
