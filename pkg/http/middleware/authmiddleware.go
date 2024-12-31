package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/brandonbraner/maas/external/usersapi"
	"github.com/brandonbraner/maas/pkg/contextservice"
	"github.com/brandonbraner/maas/pkg/errors"
	"github.com/brandonbraner/maas/pkg/http/responses"
	"github.com/brandonbraner/maas/pkg/jwtservice"
)

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

	ctxUser := contextservice.CTXUser{
		Username:    user.Username,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Permissions: user.Permissions,
		Tokens:      user.Tokens,
	}

	r = r.WithContext(context.WithValue(r.Context(), contextservice.CtxUser, ctxUser))

	next(w, r)
}
