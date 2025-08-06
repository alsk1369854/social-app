package models

type AIModelConfigs struct {
	ChatModel AIModelConfig
}

type AIModelConfig struct {
	APIKey    string
	BaseURL   string
	ModelName string
}

type AITextStreamingCallback func(chunk []byte) error

// GenerateText structs
type AIGenerateTextRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type AIGenerateTextResponse struct {
	Content string `json:"content"`
}

// GenerateTextContentOptimizationRequest structs
type AIGenerateTextContentOptimizationRequest struct {
	Context string `json:"context" binding:"required"`
	Style   string `json:"style" binding:"required"`
}

type AIGenerateTextContentOptimizationResponse struct {
	Content string `json:"content"`
}

// GenerateTextCreatePostContentRequest structs
type AIGenerateTextCreatePostContentRequest struct {
	Topic string `json:"topic" binding:"required"`
	Style string `json:"style" binding:"required"`
}

type AIGenerateTextCreatePostContentResponse struct {
	Content string `json:"content"`
}
