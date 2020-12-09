package user

import "bitbucket.org/indoquran-api/src/handlers"

const (
	userCollName = "user"
)

var (
	db             = *handlers.MongoConfig()
	userCollection = db.Collection(userCollName)
)
