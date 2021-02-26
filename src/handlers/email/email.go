package email

import (
	"fmt"
	"net/http"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/helpers"
	"bitbucket.org/indoquran-api/src/models/email"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"jaytaylor.com/html2text"
)

// SendEmail : handler to send email
func SendEmail(c *gin.Context) {
	var (
		mailData    email.Email
		mail        helpers.Email
		secret      string = c.Request.Header.Get("secret")
		timestamp   string = c.Request.Header.Get("timestamp")
		encType     string = c.Request.Header.Get("enc_type")
		isValidHMAC bool
		err         error
	)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")

	isValidHMAC, err = helpers.HMACValidation(secret, timestamp, encType)
	if err != nil {
		handlers.DefaultResponse(c, http.StatusBadRequest, "HMAC secret tidak valid", err)
		return
	}
	if !isValidHMAC {
		handlers.DefaultResponse(c, http.StatusBadRequest, "HMAC secret tidak valid", fmt.Errorf("HMAC secret tidak valid"))
		return
	}

	if err := c.ShouldBind(&mailData); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Gagal kirim Email", err)
		return
	}

	// convert body if contain html will be stripped to text only
	// Converts HTML into text of the markdown-flavored variety
	text, err := html2text.FromString(mailData.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		panic(err)
	}

	mailData.Body = text

	if err := mail.Send(c, &mailData); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Gagal kirim Email", err)
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, fmt.Sprintf("Sukses kirim Email Dari: %s, Ke: %s", mailData.From, mailData.To), nil)
	return
}
