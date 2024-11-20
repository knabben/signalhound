package llm

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
)

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI() *OpenAI {
	key := os.Getenv("OPENAI_API_KEY")
	return &OpenAI{openai.NewClient(option.WithAPIKey(key))}
}

func (o *OpenAI) RenderText(test, errMsg string) string {
	return fmt.Sprintf("A Kubernetes test titled: %s, is encountering an error: %s, summarize the error in one paragraph.", test, errMsg)
}

// PromptRequest returns the final request from OpenAI API and returns
// the LLM response to the user.
func (o *OpenAI) PromptRequest(prompt string) (string, error) {
	chatCompletion, err := o.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
