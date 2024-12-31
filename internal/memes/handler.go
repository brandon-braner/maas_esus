package memes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brandonbraner/maas/pkg/contextservice"
	"github.com/brandonbraner/maas/pkg/errors"
	"github.com/brandonbraner/maas/pkg/http/responses"
)

var memeService *MemeService

func init() {
	var err error
	memeService, err = NewMemeService()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize meme service: %v", err))
	}
}

func MemeGeneraterHandler(w http.ResponseWriter, r *http.Request) {
	var memeRequest MemeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&memeRequest)

	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}

	userctx := r.Context().Value(contextservice.CtxUser)
	ctxUser, ok := userctx.(contextservice.CTXUser)

	if !ok {
		errmsg := errors.CustomError{
			ErrorMessage: "Internal Server Error loading MemeCtx",
		}
		responses.JsonResponse(w, http.StatusInternalServerError, errmsg)
		return
	}

	aipermission := ctxUser.Permissions.GenerateLlmMeme
	memeresponse, err := memeService.GenerateMeme(aipermission, memeRequest)
	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}
	json.NewEncoder(w).Encode(memeresponse)
}
