package config

import (
	"github.com/gin-contrib/sessions/cookie"
)

// SessionWithCookies : use cookies as sesssion storage
func SessionWithCookies() cookie.Store {
	return cookie.NewStore([]byte(Config.Session.Secret))
}
