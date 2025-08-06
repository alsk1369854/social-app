export interface Post {
  id: string;
  content: string;
  userID: string;
  username?: string;
  createdAt: string;
  updatedAt?: string;
  tags?: PostTag[];
}

export interface PostTag {
  id: string;
  name: string;
}

export interface PostCreateRequest {
  content: string;
  tags?: string[];
}

export interface Comment {
  id: string;
  content: string;
  postID: string;
  userID: string;
  username?: string;
  parentID?: string;
  createdAt: string;
  updatedAt?: string;
  subComments?: Comment[];
}

export interface CommentCreateRequest {
  content: string;
  postID: string;
  parentID?: string;
}