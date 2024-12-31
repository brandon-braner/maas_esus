package memes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brandonbraner/maas/pkg/contextservice"
	"github.com/brandonbraner/maas/pkg/errors"
	"github.com/brandonbraner/maas/pkg/http/responses"
	"github.com/brandonbraner/maas/pkg/logger"
)

var log = logger.NewJsonLogger()

var memeService *MemeService

func init() {
	var err error
	memeService, err = NewMemeService()
	if err != nil {
		log.Error("Failed to initialize meme service", "error", err)
		panic(fmt.Sprintf("failed to initialize meme service: %v", err))
	}
	log.Info("Meme service initialized successfully")
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
		log.Error(err.Error())
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}

	userctx := r.Context().Value(contextservice.CtxUser)
	ctxUser, ok := userctx.(contextservice.CTXUser)

	if !ok {
		msg := "Internal Server Error loading MemeCtx"
		errmsg := errors.CustomError{
			ErrorMessage: msg,
		}
		log.Error(msg)
		responses.JsonResponse(w, http.StatusInternalServerError, errmsg)
		return
	}

	aipermission := ctxUser.Permissions.GenerateLlmMeme

	ok = memeService.VerifyTokens(aipermission, ctxUser.Tokens)
	if !ok {
		msg := fmt.Sprintf("Not enough tokens to complete request. Current token count of %d.", ctxUser.Tokens)
		errmsg := errors.CustomError{
			ErrorMessage: msg,
		}
		log.Error(msg)
		responses.JsonResponse(w, http.StatusPaymentRequired, errmsg)
		return
	}

	memeresponse, err := memeService.GenerateMeme(aipermission, memeRequest)
	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		log.Error(err.Error())
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
		msg := "Internal Server Error loading MemeCtx"
		errmsg := errors.CustomError{
			ErrorMessage: msg,
		}
		log.Error(msg)
		responses.JsonResponse(w, http.StatusInternalServerError, errmsg)
		return
	}

	tokencount, err := memeService.GetTokenCount(ctxUser.Username)

	if err != nil {
		msg := err.Error()
		errmsg := errors.CustomError{
			ErrorMessage: msg,
		}
		log.Error(err.Error())
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}

	responses.JsonResponse(w, http.StatusOK, tokencount)
}
