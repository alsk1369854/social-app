import {
  PostSearchParams,
  PostSearchResponse,
  PostCreateRequest,
  PostCreateResponse,
  ErrorResponse
} from './models/post';

const API_BASE_URL = '';

class PostAPI {
  static async createPost(request: PostCreateRequest, accessToken: string): Promise<PostCreateResponse> {
    console.log('Making POST request to /api/post with token:', accessToken);
    const response = await fetch(`${API_BASE_URL}/api/post`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': accessToken,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.error('API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to create post');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async searchPosts(params: PostSearchParams): Promise<PostSearchResponse> {
    const urlParams = new URLSearchParams();
    
    if (params.keyword) urlParams.append('keyword', params.keyword);
    if (params.offset) urlParams.append('offset', params.offset);
    if (params.limit) urlParams.append('limit', params.limit);
    if (params['user-id']) urlParams.append('user-id', params['user-id']);

    const url = `/api/post/list/search${urlParams.toString() ? `?${urlParams.toString()}` : ''}`;
    
    const response = await fetch(`${API_BASE_URL}${url}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json();
      throw new Error(errorData.error || 'Failed to search posts');
    }

    return response.json();
  }
}

export default PostAPI;