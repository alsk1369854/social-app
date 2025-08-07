export interface CommentCreateRequest {
  content: string;
  postID: string;
  parentID?: string;
}

export interface CommentCreateResponse {
  id: string;
  content: string;
  postID: string;
  parentID: string;
  userID: string;
}

export interface CommentGetListByPostIDResponseItem {
  id: string;
  content: string;
  postID: string;
  parentID: string;
  userID: string;
  userName: string;
  createdAt: string;
  updatedAt: string;
  subComments: CommentGetListByPostIDResponseItem[];
}

export interface ErrorResponse {
  error: string;
}