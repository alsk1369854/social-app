import {
  AIGenerateTextCreatePostContentRequest,
  AIGenerateTextCreatePostContentResponse,
  AIGenerateTextContentOptimizationRequest,
  AIGenerateTextContentOptimizationResponse,
  ErrorResponse
} from './models/ai';

const API_BASE_URL = '';

class AIAPI {
  static async createPostContent(request: AIGenerateTextCreatePostContentRequest, accessToken: string): Promise<AIGenerateTextCreatePostContentResponse> {
    console.log('Making POST request to /api/ai/generate/text/create-post-content with token:', accessToken);
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/create-post-content`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.error('AI API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to create AI content');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async createPostContentStream(
    request: AIGenerateTextCreatePostContentRequest, 
    accessToken: string,
    onChunk: (chunk: string) => void
  ): Promise<void> {
    console.log('Making POST request to /api/ai/generate/text/create-post-content/stream with token:', accessToken);
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/create-post-content/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.error('AI API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to create AI content stream');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('Failed to get response reader');
    }

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        
        const chunk = decoder.decode(value, { stream: true });
        onChunk(chunk);
      }
    } finally {
      reader.releaseLock();
    }
  }

  static async optimizeContent(request: AIGenerateTextContentOptimizationRequest, accessToken: string): Promise<AIGenerateTextContentOptimizationResponse> {
    console.log('Making POST request to /api/ai/generate/text/content-optimize with token:', accessToken);
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/content-optimize`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.error('AI API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to optimize content');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async optimizeContentStream(
    request: AIGenerateTextContentOptimizationRequest, 
    accessToken: string,
    onChunk: (chunk: string) => void
  ): Promise<void> {
    console.log('Making POST request to /api/ai/generate/text/content-optimize/stream with token:', accessToken);
    const response = await fetch(`${API_BASE_URL}/api/ai/generate/text/content-optimize/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.error('AI API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to optimize content stream');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('Failed to get response reader');
    }

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        
        const chunk = decoder.decode(value, { stream: true });
        onChunk(chunk);
      }
    } finally {
      reader.releaseLock();
    }
  }
}

export default AIAPI;