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

    if (!response.ok) {
      try {
        const errorData: ErrorResponse = await response.json();
        throw new Error(errorData.error || 'Failed to create AI content stream');
      } catch (jsonError) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('Failed to get response reader');
    }

    try {
      let buffer = '';
      let isCompleted = false;
      
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        
        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;
        
        let lines = buffer.split('\n');
        buffer = lines.pop() || '';
        
        for (const line of lines) {
          if (line.trim() === 'event: [DONE]') {
            isCompleted = true;
            if (onComplete) onComplete();
            return;
          }
          
          if (line.trim() === 'event: [ERROR]') {
            continue;
          }
          
          if (line.startsWith('data: ')) {
            const content = line.substring(6);
            
            if (content && content !== '' && !content.includes('[ERROR]')) {
              // Convert escaped newlines to actual newlines
              const formattedContent = content.replace(/\\n/g, '\n');
              onChunk(formattedContent);
            }
          }
        }
      }
      
      if (!isCompleted && onComplete) {
        onComplete();
      }
      
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

    if (!response.ok) {
      try {
        const errorData: ErrorResponse = await response.json();
        throw new Error(errorData.error || 'Failed to optimize content stream');
      } catch (jsonError) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('Failed to get response reader');
    }

    try {
      let buffer = '';
      let isCompleted = false;
      
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        
        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;
        
        let lines = buffer.split('\n');
        buffer = lines.pop() || '';
        
        for (const line of lines) {
          if (line.trim() === 'event: [DONE]') {
            isCompleted = true;
            if (onComplete) onComplete();
            return;
          }
          
          if (line.trim() === 'event: [ERROR]') {
            continue;
          }
          
          if (line.startsWith('data: ')) {
            const content = line.substring(6);
            
            if (content && content !== '' && !content.includes('[ERROR]')) {
              // Convert escaped newlines to actual newlines
              const formattedContent = content.replace(/\\n/g, '\n');
              onChunk(formattedContent);
            }
          }
        }
      }
      
      if (!isCompleted && onComplete) {
        onComplete();
      }
      
    } finally {
      reader.releaseLock();
    }
  }
}

export default AIAPI;