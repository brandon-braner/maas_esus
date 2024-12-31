package memes

import (
	"testing"

	"github.com/brandonbraner/maas/config"
)

func Setup(t *testing.T) *MemeService {
	config.SetupMongoTestConfig()
	service, err := NewMemeService()
	if err != nil {
		t.Fatalf("Failed to create meme service: %v", err)
	}

	return service

}

func TestMemeService_TextGenerator(t *testing.T) {

	service := Setup(t)

	req := MemeRequest{
		Query: "test",
	}

	resp, err := service.TextGenerator.Generate(req)
	if err != nil {
		t.Errorf("TextGenerator.Generate failed: %v", err)
	}

	expectedText := "This is a text meme about test"
	if resp.Text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, resp.Text)
	}
}

func TestMemeService_AITextGenerator(t *testing.T) {
	service, err := NewMemeService()
	if err != nil {
		t.Fatalf("Failed to create meme service: %v", err)
	}

	req := MemeRequest{
		Query: "test",
	}

	resp, err := service.AITextGenerator.Generate(req)
	if err != nil {
		t.Errorf("AITextGenerator.Generate failed: %v", err)
	}

	expectedText := "This is an AI-generated meme about test"
	if resp.Text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, resp.Text)
	}
}
