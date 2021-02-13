package email

import (
	"fmt"
	"net/http"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/helpers"
	"bitbucket.org/indoquran-api/src/models/email"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
)

// SendEmail : handler to send email
func SendEmail(c *gin.Context) {
	m := &email.Email{}
	mail := helpers.Email{}

	if err := c.ShouldBind(m); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Failed Send Email", err)
		return
	}

	if err := mail.Send(c, m); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusBadRequest, "Failed Send Email", err)
		return
	}

	handlers.DefaultResponse(c, http.StatusOK, fmt.Sprintf("Success Send Email From: %s, To: %s", m.From, m.To), nil)
	return
}
