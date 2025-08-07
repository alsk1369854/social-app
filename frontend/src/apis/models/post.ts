export interface ErrorResponse {
  error: string;
}

export interface Pagination {
  limit: number;
  offset: number;
}

export interface PaginationResponse<T> {
  data: T[];
  pagination: Pagination;
  totalCount: number;
}

export interface PostSearchParams {
  keyword?: string;
  offset?: string;
  limit?: string;
  'user-id'?: string;
}

export interface PostGetPostsByKeywordResponseItemTag {
  id: string;
  name: string;
}

export interface PostGetPostsByKeywordResponseItem {
  authorID: string;
  content: string;
  createdAt: string;
  id: string;
  imageURL: string;
  likedCount: number;
  tags: PostGetPostsByKeywordResponseItemTag[];
  updatedAt: string;
}

export type PostSearchResponse = PaginationResponse<PostGetPostsByKeywordResponseItem>;

export interface PostCreateRequest {
  content: string;
  imageURL?: string;
}

export interface PostCreateResponse {
  id: string;
  authorID: string;
  content: string;
  createdAt: string;
  updatedAt: string;
  imageURL: string;
  tagIDs: string[];
}