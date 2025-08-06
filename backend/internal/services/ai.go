package services

import (
	"backend/internal/models"
	"backend/internal/pkg"
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

type AIService struct {
	ErrorUtils *pkg.ErrorUtils
}

var aiServiceOnce sync.Once
var aiService *AIService

func NewAIService() *AIService {
	aiServiceOnce.Do(func() {
		aiService = &AIService{
			ErrorUtils: pkg.NewErrorUtils(),
		}
	})
	return aiService
}

func (s *AIService) CreatePostContent(ctx *gin.Context, model *openai.LLM, topic string, style string, options ...models.AITextStreamingCallback) (string, error) {
	// 使用 LangChain 的 PromptTemplate 來格式化指令
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewHumanMessagePromptTemplate(
			"你是一個專業的社群博客，請採用 {{.style}} 風格並根據使用者所提出的主題撰寫一篇博客內容。直接回覆優化後的內容，不要說任何多餘的話。",
			[]string{"style"},
		),
		prompts.NewHumanMessagePromptTemplate(
			`請撰寫一篇關於: {{.topic}} 的博客內容`,
			[]string{"content"},
		),
	})
	instruction, err := prompt.Format(map[string]any{
		"topic": topic,
		"style": style,
	})
	if err != nil {
		return "", err
	}

	// 使用 LangChain 的 LLM 生成內容
	output, err := llms.GenerateFromSinglePrompt(
		context.Background(), model, instruction,
		llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
			for _, callback := range options {
				if callback == nil {
					continue
				}
				if err := callback(chunk); err != nil {
					return err
				}
			}
			return nil
		}),
	)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (s *AIService) ContentOptimization(ctx *gin.Context, model *openai.LLM, content string, style string, options ...models.AITextStreamingCallback) (string, error) {

	// 使用 LangChain 的 PromptTemplate 來格式化指令
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			"你是一個專業的社群博客，請優化使用者輸入的內容以符合 {{.style}} 的風格。直接回覆優化後的內容，不要說任何多餘的話。",
			[]string{"style"},
		),
		prompts.NewHumanMessagePromptTemplate(
			`{{.content}}`,
			[]string{"content"},
		),
	})
	instruction, err := prompt.Format(map[string]any{
		"content": content,
		"style":   style,
	})
	if err != nil {
		return "", err
	}

	// 使用 LangChain 的 LLM 生成內容
	output, err := llms.GenerateFromSinglePrompt(
		context.Background(), model, instruction,
		llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
			for _, callback := range options {
				if callback == nil {
					continue
				}
				if err := callback(chunk); err != nil {
					return err
				}
			}
			return nil
		}),
	)
	if err != nil {
		return "", err
	}

	return output, nil
}
