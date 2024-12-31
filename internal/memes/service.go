package memes

import (
	"context"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/external/usersapi"
)

type MemeService struct {
	Repo            *memeRepository
	TextGenerator   *MemeGenerator // Strategy for text memes
	AITextGenerator *MemeGenerator // Strategy for AI memes
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
	service := &MemeService{
		Repo:            repo,
		TextGenerator:   textgen,
		AITextGenerator: aitextgen,
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
	tokenservice, err := usersapi.NewUserService()
	if err != nil {
		return err
	}

	var numOfTokens int
	if aiGenerated {
		numOfTokens = config.AppConfig.AI_TEXT_MEME_TOKEN_COST
	} else {
		numOfTokens = config.AppConfig.TEXT_MEME_TOKEN_COST
	}
	//turn tokens negative
	numOfTokens = numOfTokens * -1

	tokenservice.UpdateTokens(username, numOfTokens)
	return nil
}
