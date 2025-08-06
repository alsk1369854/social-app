import React from 'react';
import Post from './Post';

interface Comment {
  id: string;
  content: string;
  userID: string;
  username?: string;
  createdAt: string;
  subComments?: Comment[];
}

interface PostTag {
  id: string;
  name: string;
}

interface PostData {
  id: string;
  content: string;
  userID: string;
  username?: string;
  createdAt: string;
  tags?: PostTag[];
}

interface PostsFeedProps {
  posts: PostData[];
  postComments: Record<string, Comment[]>;
  isLoggedIn: boolean;
  loading?: boolean;
  onAddComment: (postId: string, content: string) => Promise<void>;
  onLoadComments: (postId: string) => Promise<void>;
}

const PostsFeed: React.FC<PostsFeedProps> = ({
  posts,
  postComments,
  isLoggedIn,
  loading = false,
  onAddComment,
  onLoadComments
}) => {
  if (loading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3].map((index) => (
          <div key={index} className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 animate-pulse">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-10 h-10 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
              <div>
                <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-24 mb-2"></div>
                <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-16"></div>
              </div>
            </div>
            <div className="space-y-2 mb-4">
              <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-full"></div>
              <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-3/4"></div>
              <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/2"></div>
            </div>
            <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
              <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-20"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (posts.length === 0) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-8 text-center">
        <div className="text-gray-500 dark:text-gray-400 mb-4">
          <svg
            className="w-12 h-12 mx-auto mb-4 text-gray-300 dark:text-gray-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            還沒有貼文
          </h3>
          <p className="text-gray-500 dark:text-gray-400">
            成為第一個分享想法的人吧！
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <Post
          key={post.id}
          post={post}
          comments={postComments[post.id] || []}
          isLoggedIn={isLoggedIn}
          onAddComment={onAddComment}
          onLoadComments={onLoadComments}
        />
      ))}
    </div>
  );
};

export default PostsFeed;