export interface ErrorResponse {
  error: string;
}

export interface AIGenerateTextCreatePostContentRequest {
  topic: string;
  style: string;
}

export interface AIGenerateTextCreatePostContentResponse {
  content: string;
}

export interface AIGenerateTextContentOptimizationRequest {
  context: string;
  style: string;
}

export interface AIGenerateTextContentOptimizationResponse {
  content: string;
}