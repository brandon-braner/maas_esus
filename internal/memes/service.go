package memes

import "context"

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
