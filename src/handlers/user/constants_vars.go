package user

import "bitbucket.org/indoquran-api/src/handlers"

const (
	userCollName         = "user"
	loginHistoryCollName = "login_history"
)

var (
	db                     = *handlers.MongoInstance()
	userCollection         = db.Collection(userCollName)
	loginHistoryCollection = db.Collection(loginHistoryCollName)
)
