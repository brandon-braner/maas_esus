package memes

import (
	"context"
	"fmt"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/external/usersapi"
	"github.com/brandonbraner/maas/pkg/contextservice"
)

type MemeService struct {
	Repo            *memeRepository
	TextGenerator   *MemeGenerator // Strategy for text memes
	AITextGenerator *MemeGenerator // Strategy for AI memes
	UserService     *usersapi.UserService
}

func NewMemeService() (*MemeService, error) {
	repo, err := NewMemeRepository(context.Background())

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	textgen, err := NewMemeGenerator(false)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	aitextgen, err := NewMemeGenerator(true)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	userservice, err := usersapi.NewUserService()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	service := &MemeService{
		Repo:            repo,
		TextGenerator:   textgen,
		AITextGenerator: aitextgen,
		UserService:     userservice,
	}

	return service, nil
}

func (s *MemeService) GenerateMeme(aiPermission bool, memeRequest MemeRequest) (MemeResponse, error) {
	switch aiPermission {
	case true:
		return (*s.AITextGenerator).Generate(memeRequest)

	default:
		return (*s.TextGenerator).Generate(memeRequest)
	}

}

func (s *MemeService) ChargeTokens(aiGenerated bool, username string) error {

	var numOfTokens int
	if aiGenerated {
		numOfTokens = config.AppConfig.AI_TEXT_MEME_TOKEN_COST
	} else {
		numOfTokens = config.AppConfig.TEXT_MEME_TOKEN_COST
	}
	//turn tokens negative
	tokensToCharge := numOfTokens * -1

	err := s.UserService.UpdateTokens(username, tokensToCharge)
	if err != nil {
		log.Error(fmt.Sprintf("Could not charge user %s token amount %d. Still returning meme", username, numOfTokens))
	}
	log.Info(fmt.Sprintf("%d tokens charged to %s", numOfTokens, username))
	return nil
}
func (s *MemeService) VerifyTokens(aiGenerated bool, user contextservice.CTXUser) bool {

	if user.Tokens < 0 {
		return false
	}

	var tokensRequired int
	if aiGenerated {
		tokensRequired = config.AppConfig.AI_TEXT_MEME_TOKEN_COST
	} else {
		tokensRequired = config.AppConfig.TEXT_MEME_TOKEN_COST
	}

	if user.Tokens < tokensRequired {
		return false
	}
	return true

}

func (s *MemeService) GetTokenCount(username string) (int, error) {

	tokencount, err := s.UserService.GetTokenCount(username)

	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	return tokencount, nil
}
