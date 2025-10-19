import { on } from 'events';
import {
  AIGenerateTextCreatePostContentRequest,
  AIGenerateTextCreatePostContentResponse,
  AIGenerateTextContentOptimizationRequest,
  AIGenerateTextContentOptimizationResponse,
  ErrorResponse
} from './models/ai';

// const API_BASE_URL = 'http://localhost:28080';
const API_BASE_URL = process.env.REACT_APP_BASE_URL || "";


class AIAPI {
  static async createPostContent(request: AIGenerateTextCreatePostContentRequest, accessToken: string): Promise<AIGenerateTextCreatePostContentResponse> {
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/create-post-content`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      try {
        const errorData: ErrorResponse = await response.json();
        throw new Error(errorData.error || 'Failed to create AI content');
      } catch (jsonError) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async createPostContentStream(
    request: AIGenerateTextCreatePostContentRequest,
    accessToken: string,
    onChunk: (chunk: string) => void,
    onComplete?: () => void,
    onError?: (error: string) => void
  ): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/create-post-content/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok || !response.body) {
      if (onError) { onError("sse connection failed"); }
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        const event = decoder.decode(value, { stream: true });
        const chunks = event.split("data: ")
        for (const chunk of chunks) {
          const data = chunk.substring(0, chunk.length - 2)
          if (data === "[DONE]") {
            if (onComplete) onComplete()
            return
          }
          if (data.startsWith("[ERROR]")) {
            if (onError) onError("server busy please try again later");
            if (onComplete) onComplete()
            return
          }
          onChunk(data)
        }
      }
    } catch (error) {
      if (onError) onError((error as Error).message);
    } finally {
      reader.releaseLock();
    }
  }

  static async optimizeContent(request: AIGenerateTextContentOptimizationRequest, accessToken: string): Promise<AIGenerateTextContentOptimizationResponse> {
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/content-optimize`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      try {
        const errorData: ErrorResponse = await response.json();
        throw new Error(errorData.error || 'Failed to optimize content');
      } catch (jsonError) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async optimizeContentStream(
    request: AIGenerateTextContentOptimizationRequest,
    accessToken: string,
    onChunk: (chunk: string) => void,
    onComplete?: () => void,
    onError?: (error: string) => void
  ): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/content-optimize/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok || !response.body) {
      if (onError) { onError("sse connection failed"); }
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        const event = decoder.decode(value, { stream: true });
        const chunks = event.split("data: ")
        for (const chunk of chunks) {
          const data = chunk.substring(0, chunk.length - 2)
          if (data === "[DONE]") {
            if (onComplete) onComplete()
            return
          }
          if (data.startsWith("[ERROR]")) {
            if (onError) onError("server busy please try again later");
            if (onComplete) onComplete()
            return
          }
          onChunk(data)
        }
      }
    } catch (error) {
      if (onError) onError((error as Error).message);
    } finally {
      reader.releaseLock();
    }
  }
}

export default AIAPI;