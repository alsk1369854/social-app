import React, { useState } from 'react';

interface PostCreatorProps {
  isLoggedIn: boolean;
  onCreatePost: (content: string) => Promise<void>;
}

const PostCreator: React.FC<PostCreatorProps> = ({ isLoggedIn, onCreatePost }) => {
  const [content, setContent] = useState('');
  const [isExpanded, setIsExpanded] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const maxLength = 300;
  const remainingChars = maxLength - content.length;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim() || !isLoggedIn || isSubmitting) return;

    setIsSubmitting(true);
    try {
      await onCreatePost(content.trim());
      setContent('');
      setIsExpanded(false);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleTextareaClick = () => {
    if (isLoggedIn) {
      setIsExpanded(true);
    }
  };

  const handleCancel = () => {
    setContent('');
    setIsExpanded(false);
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 mb-6">
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            onClick={handleTextareaClick}
            disabled={!isLoggedIn}
            maxLength={maxLength}
            className={`w-full p-3 sm:p-4 border border-gray-300 dark:border-gray-600 rounded-md resize-none bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all ${
              !isLoggedIn ? 'cursor-not-allowed opacity-50' : ''
            } ${isExpanded ? 'h-32 sm:h-36' : 'h-12 sm:h-14'}`}
            placeholder={isLoggedIn ? "分享你的想法..." : "請先登入以發布貼文"}
          />
        </div>

        {isExpanded && (
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-2 sm:space-y-0">
            <div className="flex items-center">
              <span className={`text-sm ${
                remainingChars < 0 
                  ? 'text-red-500' 
                  : remainingChars < 50 
                    ? 'text-yellow-500' 
                    : 'text-gray-500 dark:text-gray-400'
              }`}>
                {remainingChars} 字元剩餘
              </span>
            </div>
            
            <div className="flex space-x-2 w-full sm:w-auto justify-end">
              <button
                type="button"
                onClick={handleCancel}
                disabled={isSubmitting}
                className="px-3 sm:px-4 py-2 text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 text-sm font-medium transition-colors disabled:opacity-50"
              >
                取消
              </button>
              <button
                type="submit"
                disabled={!content.trim() || remainingChars < 0 || isSubmitting}
                className="px-3 sm:px-4 py-2 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white rounded-md text-sm font-medium transition-colors"
              >
                {isSubmitting ? '發布中...' : '發布'}
              </button>
            </div>
          </div>
        )}
      </form>
    </div>
  );
};

export default PostCreator;