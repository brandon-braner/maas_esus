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
	// var err error

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

	ok = memeService.VerifyTokens(aipermission, ctxUser.Tokens)
	if !ok {
		errmsg := errors.CustomError{
			ErrorMessage: fmt.Sprintf("Not enough tokens to complete request. Current token count of %d.", ctxUser.Tokens),
		}
		responses.JsonResponse(w, http.StatusPaymentRequired, errmsg)
		return
	}

	memeresponse, err := memeService.GenerateMeme(aipermission, memeRequest)
	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}

	//Assume we have made it here, we got a meme, lets charge them some tokens
	memeService.ChargeTokens(aipermission, ctxUser.Username)
	json.NewEncoder(w).Encode(memeresponse)
}

func MemeTokenHandler(w http.ResponseWriter, r *http.Request) {
	userctx := r.Context().Value(contextservice.CtxUser)
	ctxUser, ok := userctx.(contextservice.CTXUser)

	if !ok {
		errmsg := errors.CustomError{
			ErrorMessage: "Internal Server Error loading MemeCtx",
		}
		responses.JsonResponse(w, http.StatusInternalServerError, errmsg)
		return
	}

	tokencount, err := memeService.GetTokenCount(ctxUser.Username)

	if err != nil {
		if err != nil {
			errmsg := errors.CustomError{
				ErrorMessage: err.Error(),
			}
			responses.JsonResponse(w, http.StatusBadRequest, errmsg)
			return
		}
	}

	responses.JsonResponse(w, http.StatusOK, tokencount)
}
