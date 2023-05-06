package chat

import (
	"chatbot/common"
	"context"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const AIModel = "text-davinci-003"

type GptResponse struct {
	Content          string
	PromptToken      int
	CompletionTokens int
}

type Service interface {
	GoGpt(question string, info string) GptResponse
}

type service struct {
	gptClient *gogpt.Client
}

func NewService(key string) Service {
	client := gogpt.NewClient(key)

	return &service{
		gptClient: client,
	}
}

func (s *service) GoGpt(question string, info string) GptResponse {

	req := gogpt.CompletionRequest{
		Model:       AIModel,
		MaxTokens:   100,
		Prompt:      info + question + ".",
		Temperature: 0.1,
	}
	common.Sugar.Infow("Asking the Davinci", "question", info+question)

	resp, err := s.gptClient.CreateCompletion(context.Background(), req)

	if err != nil {
		common.Sugar.Fatalw("Error on getting the response from the Davinci API", "error", err)
	}

	common.Sugar.Infow("Successfully got the response from the Davinci API", "response", resp)

	pt := resp.Usage.PromptTokens

	ct := resp.Usage.CompletionTokens

	content := resp.Choices[0].Text

	return GptResponse{
		Content:          content,
		PromptToken:      pt,
		CompletionTokens: ct,
	}

}
