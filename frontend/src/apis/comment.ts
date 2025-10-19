import {
  CommentCreateRequest,
  CommentCreateResponse,
  CommentGetListByPostIDResponseItem,
  ErrorResponse
} from './models/comment';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || "";

class CommentAPI {
  static async createComment(request: CommentCreateRequest, accessToken: string): Promise<CommentCreateResponse> {
    const response = await fetch(`${API_BASE_URL}/comment`, {
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
        throw new Error(errorData.error || 'Failed to create comment');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }

  static async getCommentsByPostID(postID: string): Promise<CommentGetListByPostIDResponseItem[]> {
    const response = await fetch(`${API_BASE_URL}/comment/list/post/${postID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      console.error('API Error Response:', response.status, response.statusText);
      try {
        const errorData: ErrorResponse = await response.json();
        console.error('Error Data:', errorData);
        throw new Error(errorData.error || 'Failed to load comments');
      } catch (jsonError) {
        console.error('Failed to parse error response:', jsonError);
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
    }

    return response.json();
  }
}

export default CommentAPI;