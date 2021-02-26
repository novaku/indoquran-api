package helpers

import (
	"crypto/tls"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/models/email"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// Email : email structure helpers
type Email struct{}

// Send : send email
func (e *Email) Send(ctx *gin.Context, m *email.Email) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.From)
	mail.SetHeader("To", m.To)
	mail.SetAddressHeader("Cc", m.From, m.From)
	mail.SetHeader("Subject", m.Subject)
	mail.SetBody("text/plain", m.Body)

	d := gomail.NewDialer(config.Config.Email.SMTP,
		config.Config.Email.Port,
		config.Config.Email.User,
		config.Config.Email.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
