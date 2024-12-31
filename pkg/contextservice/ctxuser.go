package contextservice

import (
	"github.com/brandonbraner/maas/pkg/permissions"
)

type contextKey string

const CtxUser contextKey = "ctxuser"

type CTXUser struct {
	Username                string `json:"username" bson:"username"`
	Firstname               string `json:"firstname" bson:"firstname"`
	Lastname                string `json:"lastname" bson:"lastname"`
	permissions.Permissions `json:"permissions" bson:"permissions"`
}
