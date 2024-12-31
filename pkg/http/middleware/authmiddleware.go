package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/brandonbraner/maas/external/usersapi"
	"github.com/brandonbraner/maas/pkg/errors"
	"github.com/brandonbraner/maas/pkg/http/responses"
	"github.com/brandonbraner/maas/pkg/jwtservice"
)

type contextKey string

const CtxUser contextKey = "ctxuser"

type AuthMiddleware struct{}

func (am AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := errors.CustomError{
			ErrorMessage: "Authorization Header Not Found",
		}

		responses.JsonResponse(w, http.StatusUnauthorized, err)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwtservice.ValidateToken(token)

	if err != nil {
		err := errors.CustomError{
			ErrorMessage: err.Error(),
		}

		responses.JsonResponse(w, http.StatusUnauthorized, err)
		return
	}

	// mongoClient, err := database.GetClient()
	// if err != nil {
	// w.WriteHeader(http.StatusInternalServerError)
	// return
	// }
	userService, err := usersapi.NewUserService()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := userService.GetUserByUsername(claims.Email)
	if err != nil {
		customerr := errors.CustomError{
			ErrorMessage: "User not found",
		}
		responses.JsonResponse(w, http.StatusUnauthorized, customerr)
		return
	}

	// ctxPermissions, err := contextservice.ConvertViaJSON(u.Permissions)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// ctxUser := contextservice.User{
	// 	Username:    u.Username,
	// 	Firstname:   u.Firstname,
	// 	Lastname:    u.Lastname,
	// 	Permissions: ctxPermissions,
	// }

	r = r.WithContext(context.WithValue(r.Context(), CtxUser, user))

	next(w, r)
}
