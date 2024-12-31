package memes

import (
	"context"
	"fmt"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/external/usersapi"
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
		return nil, err
	}

	textgen, err := NewMemeGenerator(false)
	if err != nil {
		return nil, err
	}

	aitextgen, err := NewMemeGenerator(true)
	if err != nil {
		return nil, err
	}

	userservice, err := usersapi.NewUserService()
	if err != nil {
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
	numOfTokens = numOfTokens * -1

	err := s.UserService.UpdateTokens(username, numOfTokens)
	if err != nil {
		fmt.Sprintf("Could not charge user %s token amount %d. Still returning meme", username, numOfTokens)
	}
	return nil
}

func (s *MemeService) VerifyTokens(aiGenerated bool, currenttokens int) bool {
	if currenttokens < 0 {
		return false
	}

	var tokensRequired int
	if aiGenerated {
		tokensRequired = config.AppConfig.AI_TEXT_MEME_TOKEN_COST
	} else {
		tokensRequired = config.AppConfig.TEXT_MEME_TOKEN_COST
	}

	if currenttokens < tokensRequired {
		return false
	}
	return true

}

func (s *MemeService) GetTokenCount(username string) (int, error) {

	tokencount, err := s.UserService.GetTokenCount(username)

	if err != nil {
		return 0, err
	}

	return tokencount, nil
}
