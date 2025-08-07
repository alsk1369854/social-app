import React, { useState } from 'react';
import MarkdownRenderer from './MarkdownRenderer';
import MarkdownEditor from './MarkdownEditor';

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

interface PostProps {
  post: PostData;
  comments?: Comment[];
  isLoggedIn: boolean;
  onAddComment: (postId: string, content: string) => Promise<void>;
  onLoadComments: (postId: string) => Promise<void>;
}

const Post: React.FC<PostProps> = ({
  post,
  comments = [],
  isLoggedIn,
  onAddComment,
  onLoadComments
}) => {
  const [showComments, setShowComments] = useState(false);
  const [commentContent, setCommentContent] = useState('');
  const [isSubmittingComment, setIsSubmittingComment] = useState(false);
  const [commentsLoaded, setCommentsLoaded] = useState(false);

  const handleToggleComments = async () => {
    if (!showComments && !commentsLoaded) {
      await onLoadComments(post.id);
      setCommentsLoaded(true);
    }
    setShowComments(!showComments);
  };

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentContent.trim() || !isLoggedIn || isSubmittingComment) return;

    setIsSubmittingComment(true);
    try {
      await onAddComment(post.id, commentContent.trim());
      setCommentContent('');
    } finally {
      setIsSubmittingComment(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-TW', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const renderComment = (comment: Comment, isSubComment = false) => (
    <div key={comment.id} className={`${isSubComment ? 'ml-4 sm:ml-6 pl-2 sm:pl-4 border-l border-gray-200 dark:border-gray-600' : ''} py-2`}>
      <div className="flex items-start space-x-2 sm:space-x-3">
        <div className="w-6 h-6 sm:w-8 sm:h-8 bg-gray-500 rounded-full flex items-center justify-center text-white font-semibold text-xs">
          {comment.username?.charAt(0).toUpperCase() || 'U'}
        </div>
        <div className="flex-1 min-w-0">
          <div className="bg-gray-100 dark:bg-gray-700 rounded-lg p-2 sm:p-3">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-1">
              <span className="font-medium text-sm text-gray-900 dark:text-white truncate">
                {comment.username || 'Unknown User'}
              </span>
              <span className="text-xs text-gray-500 dark:text-gray-400 flex-shrink-0">
                {formatDate(comment.createdAt)}
              </span>
            </div>
            <MarkdownRenderer 
              content={comment.content} 
              className="text-sm text-gray-800 dark:text-gray-200" 
            />
          </div>
          {comment.subComments && comment.subComments.length > 0 && (
            <div className="mt-2">
              {comment.subComments.map(subComment => renderComment(subComment, true))}
            </div>
          )}
        </div>
      </div>
    </div>
  );

  return (
    <article className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 sm:p-6 mb-4">
      {/* Post Header */}
      <div className="flex items-center space-x-3 mb-4">
        <div className="w-8 h-8 sm:w-10 sm:h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-semibold text-sm">
          {post.username?.charAt(0).toUpperCase() || 'U'}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-medium text-gray-900 dark:text-white text-sm sm:text-base truncate">
            {post.username || 'Unknown User'}
          </h3>
          <p className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
            {formatDate(post.createdAt)}
          </p>
        </div>
      </div>

      {/* Post Content */}
      <div className="mb-4">
        <MarkdownRenderer 
          content={post.content} 
          className="text-sm sm:text-base leading-relaxed" 
        />
      </div>

      {/* Post Tags */}
      {post.tags && post.tags.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-4">
          {post.tags.map((tag, index) => (
            <span
              key={index}
              className="inline-block bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs px-2 py-1 rounded-full"
            >
              #{tag.name}
            </span>
          ))}
        </div>
      )}

      {/* Post Actions */}
      <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
        <button
          onClick={handleToggleComments}
          className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 text-sm font-medium"
        >
          {showComments ? '隱藏留言' : `查看留言 (${comments.length})`}
        </button>
      </div>

      {/* Comments Section */}
      {showComments && (
        <div className="mt-4 border-t border-gray-200 dark:border-gray-700 pt-4">
          {/* Comments List */}
          {comments.length > 0 && (
            <div className="mb-4 space-y-2">
              {comments.map(comment => renderComment(comment))}
            </div>
          )}

          {/* Add Comment Form */}
          {isLoggedIn ? (
            <form onSubmit={handleSubmitComment} className="mt-4">
              <div className="mb-3">
                <MarkdownEditor
                  value={commentContent}
                  onChange={setCommentContent}
                  placeholder="寫個留言..."
                  className="min-h-16"
                />
              </div>
              <div className="flex justify-between items-center">
                <div className="text-xs text-gray-400">
                  {commentContent.length > 0 && `${commentContent.length} 字元`}
                </div>
                <div className="flex space-x-2">
                  <button
                    type="button"
                    onClick={() => setCommentContent('')}
                    className="px-3 py-1 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 text-sm transition-colors"
                  >
                    清除
                  </button>
                  <button
                    type="submit"
                    disabled={!commentContent.trim() || isSubmittingComment}
                    className="px-3 py-1 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white rounded-md text-sm font-medium transition-colors"
                  >
                    {isSubmittingComment ? '發布中...' : '留言'}
                  </button>
                </div>
              </div>
            </form>
          ) : (
            <div className="mt-4 p-3 sm:p-4 bg-gray-100 dark:bg-gray-700 rounded-lg text-center">
              <p className="text-gray-600 dark:text-gray-400 text-sm">
                請先登入以留言
              </p>
            </div>
          )}
        </div>
      )}
    </article>
  );
};

export default Post;