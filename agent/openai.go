package agent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"log"
	"os"
)

func createOpenAIChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: os.Getenv("OPENAI_API_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
