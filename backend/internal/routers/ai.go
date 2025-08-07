package routers

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/pkg"
	"backend/internal/services"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type AIRouter struct {
	ErrorUtils *pkg.ErrorUtils

	ChatModel *openai.LLM
	AIService *services.AIService
}

var aiRouterOnce sync.Once
var aiRouter *AIRouter

func NewAIRouter(modelsConfigs *models.AIModelConfigs) *AIRouter {
	aiRouterOnce.Do(func() {
		chatModel, err := openai.New(
			openai.WithToken(modelsConfigs.ChatModel.APIKey),
			openai.WithBaseURL(modelsConfigs.ChatModel.BaseURL),
			openai.WithModel(modelsConfigs.ChatModel.ModelName),
		)
		if err != nil {
			log.Fatal(err)
		}

		aiRouter = &AIRouter{
			ErrorUtils: pkg.NewErrorUtils(),

			ChatModel: chatModel,
			AIService: services.NewAIService(),
		}
	})
	return aiRouter
}

func (r *AIRouter) Bind(_router *gin.RouterGroup) {
	router := _router.Group("/ai")
	// POST
	{
		router.POST("/generate/text", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.GenerateText)
		router.POST("/generate/text/stream", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.GenerateTextStream)
		router.POST("/generate/text/content-optimize", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.ContentOptimization)
		router.POST("/generate/text/content-optimize/stream", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.ContentOptimizationStream)
		router.POST("/generate/text/create-post-content", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.CreatePostContent)
		router.POST("/generate/text/create-post-content/stream", middlewares.VerifyAccessToken(middlewares.ParseJWTAccessToken), r.CreatePostContentStream)
	}
}

// @title AI API
// @Summary Create post content using AI with streaming
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce text/event-stream
// @Param request body models.AIGenerateTextCreatePostContentRequest true "AI Create Post Content Request"
// @Success 200 {string} string "Streaming response"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text/create-post-content/stream [post]
func (r *AIRouter) CreatePostContentStream(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextCreatePostContentRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// 設定 Header 為流式傳輸
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	if _, err := r.AIService.CreatePostContent(
		ctx, r.ChatModel, reqBody.Topic, reqBody.Style,
		func(chunk []byte) error {
			formatted := fmt.Sprintf("data: %s\n\n", chunk)
			if _, err := ctx.Writer.Write([]byte(formatted)); err != nil {
				return err
			}
			ctx.Writer.Flush()
			return nil
		},
	); err != nil {
		fmt.Fprintf(ctx.Writer, "event: [ERROR]\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}

	// 結束訊號
	fmt.Fprintf(ctx.Writer, "event: [DONE]\n\n")
	ctx.Writer.Flush()
}

// @title AI API
// @Summary Create post content using AI
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce application/json
// @Param request body models.AIGenerateTextCreatePostContentRequest true "AI Create Post Content Request"
// @Success 200 {object} models.AIGenerateTextCreatePostContentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text/create-post-content [post]
func (r *AIRouter) CreatePostContent(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextCreatePostContentRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	output, err := r.AIService.CreatePostContent(ctx, r.ChatModel, reqBody.Topic, reqBody.Style)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: err.Error()})
		return
	}

	respBody := models.AIGenerateTextCreatePostContentResponse{
		Content: output,
	}
	ctx.JSON(200, respBody)
}

// @title AI API
// @Summary Optimize content using AI
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce application/json
// @Param request body models.AIGenerateTextContentOptimizationRequest true "AI Content Generation Request"
// @Success 200 {object} models.AIGenerateTextContentOptimizationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text/content-optimize [post]
func (r *AIRouter) ContentOptimization(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextContentOptimizationRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	output, err := r.AIService.ContentOptimization(ctx, r.ChatModel, reqBody.Context, reqBody.Style)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: "Failed to optimize content"})
		return
	}

	respBody := models.AIGenerateTextContentOptimizationResponse{
		Content: output,
	}
	ctx.JSON(200, respBody)
}

// @title AI API
// @Summary Optimize content using AI with streaming
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce text/event-stream
// @Param request body models.AIGenerateTextContentOptimizationRequest true "AI Content Optimization Request"
// @Success 200 {string} string "Streaming response"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text/content-optimize/stream [post]
func (r *AIRouter) ContentOptimizationStream(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextContentOptimizationRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// 設定 Header 為流式傳輸
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	if _, err := r.AIService.ContentOptimization(
		ctx, r.ChatModel, reqBody.Context, reqBody.Style,
		func(chunk []byte) error {
			formatted := fmt.Sprintf("data: %s\n\n", chunk)
			if _, err := ctx.Writer.Write([]byte(formatted)); err != nil {
				return err
			}
			ctx.Writer.Flush()
			return nil
		},
	); err != nil {
		fmt.Fprintf(ctx.Writer, "event: [ERROR]\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}

	// 結束訊號
	fmt.Fprintf(ctx.Writer, "event: [DONE]\n\n")
	ctx.Writer.Flush()
}

// @title AI API
// @Summary Generate content using AI with streaming
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce text/event-stream
// @Param request body models.AIGenerateTextRequest true "AI Content Generation Request"
// @Success 200 {string} string "Streaming response"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text/stream [post]
func (r *AIRouter) GenerateTextStream(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// 設定 Header 為流式傳輸
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	if _, err := llms.GenerateFromSinglePrompt(
		context.Background(), r.ChatModel, reqBody.Prompt,
		llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
			formatted := fmt.Sprintf("data: %s\n\n", chunk)
			if _, err := ctx.Writer.Write([]byte(formatted)); err != nil {
				return err
			}
			ctx.Writer.Flush()
			return nil
		}),
	); err != nil {
		fmt.Fprintf(ctx.Writer, "event: [ERROR]\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}

	// 結束訊號
	fmt.Fprintf(ctx.Writer, "event: [DONE]\n\n")
	ctx.Writer.Flush()
}

// @title AI API
// @Summary Generate content using AI
// @Tags AI
// @Security AccessToken
// @Accept application/json
// @Produce application/json
// @Param request body models.AIGenerateTextRequest true "AI Content Generation Request"
// @Success 200 {object} models.AIGenerateTextResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/ai/generate/text [post]
func (r *AIRouter) GenerateText(ctx *gin.Context) {
	reqBody := &models.AIGenerateTextRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	output, err := llms.GenerateFromSinglePrompt(
		context.Background(), r.ChatModel, reqBody.Prompt,
	)
	if err != nil {
		ctx.JSON(500, models.ErrorResponse{Error: "Failed to generate content"})
		return
	}

	respBody := models.AIGenerateTextResponse{
		Content: output,
	}
	ctx.JSON(200, respBody)
}
