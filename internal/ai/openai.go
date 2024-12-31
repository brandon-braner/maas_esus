package ai

import (
	"context"
	"github.com/openai/openai-go"
	//"github.com/openai/openai-go/option"
)

type MemePrompt struct {
	SystemPrompt string `json:"systemPrompt"`
	UserPrompt   string `json:"userPrompt"`
}

type OpenAIMemeService struct {
	client *openai.Client
}

func (ai *OpenAIMemeService) GenerateTextMeme(prompts *MemePrompt) (string, error) {
	chatCompletion, err := ai.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompts.SystemPrompt),
			openai.UserMessage(prompts.UserPrompt),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})

	if err != nil {
		panic(err.Error())
	}

	message := chatCompletion.Choices[0].Message.Content
	println(message)
	return message, nil
}

func NewOpenAIMemeService() *OpenAIMemeService {
	client := openai.NewClient() // defaults to os.LookupEnv("OPENAI_API_KEY") )
	return &OpenAIMemeService{client: client}
}
