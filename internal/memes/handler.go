package memes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"

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
	ctx := r.Context()
	tracer := otel.Tracer("memes")

	// Start span for request processing
	ctx, span := tracer.Start(ctx, "meme-generator-handler")
	defer span.End()

	// Span for request decoding
	_, decodeSpan := tracer.Start(ctx, "decode-request")
	var memeRequest MemeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&memeRequest)

	// add request context for opentelemetry to trace the whole call
	memeRequest.Context = ctx

	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		log.Error(err.Error())
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}
	decodeSpan.End()

	// Span for user context extraction
	_, userCtxSpan := tracer.Start(ctx, "extract-user-context")
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
	userCtxSpan.End()

	aipermission := ctxUser.Permissions.GenerateLlmMeme

	// Span for token verification
	_, verifySpan := tracer.Start(ctx, "verify-tokens")
	ok = memeService.VerifyTokens(aipermission, ctxUser)

	if !ok {
		msg := fmt.Sprintf("Not enough tokens to complete request. Current token count of %d.", ctxUser.Tokens)
		errmsg := errors.CustomError{
			ErrorMessage: msg,
		}
		log.Error(msg)
		responses.JsonResponse(w, http.StatusPaymentRequired, errmsg)
		return
	}
	verifySpan.End()

	// Span for meme generation
	_, generateSpan := tracer.Start(ctx, "generate-meme")
	memeresponse, err := memeService.GenerateMeme(aipermission, memeRequest)

	if err != nil {
		errmsg := errors.CustomError{
			ErrorMessage: err.Error(),
		}
		log.Error(err.Error())
		responses.JsonResponse(w, http.StatusBadRequest, errmsg)
		return
	}
	generateSpan.End()

	// Span for token charging
	_, chargeSpan := tracer.Start(ctx, "charge-tokens")
	memeService.ChargeTokens(aipermission, ctxUser.Username)
	chargeSpan.End()

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
