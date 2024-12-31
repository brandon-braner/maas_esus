package memes

import (
	"testing"

	"github.com/brandonbraner/maas/config"
)

type MockMemeGenerator struct {
	called bool
}

func (m *MockMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {
	m.called = true
	return MemeResponse{Text: "mocked response"}, nil
}

func Setup(t *testing.T) *MemeService {
	config.SetupMongoTestConfig()
	service, err := NewMemeService()
	if err != nil {
		t.Fatalf("Failed to create meme service: %v", err)
	}

	return service
}

func TestMemeService_GenerateMeme_Text(t *testing.T) {
	service := Setup(t)

	req := MemeRequest{
		Query: "test",
	}

	resp, err := service.GenerateMeme(false, req)
	if err != nil {
		t.Errorf("GenerateMeme failed: %v", err)
	}

	expectedText := "This is a text meme about test"
	if resp.Text != expectedText {
		t.Errorf("Expected text %q, got %q", expectedText, resp.Text)
	}
}

func TestMemeService_GenerateMeme_AI(t *testing.T) {
	service := Setup(t)
	
	// Replace AI generator with mock
	mockAI := &MockMemeGenerator{}
	var generator MemeGenerator = mockAI
	service.AITextGenerator = &generator

	req := MemeRequest{
		Query: "test",
	}

	_, err := service.GenerateMeme(true, req)
	if err != nil {
		t.Errorf("GenerateMeme failed: %v", err)
	}

	if !mockAI.called {
		t.Error("Expected AI generator to be called but it wasn't")
	}
}

func TestMemeService_GenerateMeme_Error(t *testing.T) {
	service := &MemeService{} // Empty service to force error

	req := MemeRequest{
		Query: "test",
	}

	_, err := service.GenerateMeme(true, req)
	if err == nil {
		t.Error("Expected error from GenerateMeme with invalid service")
	}
}
