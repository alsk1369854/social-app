import React, { useState, useContext, use, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AIToolModal from './AIToolModal';
import MarkdownEditor from './MarkdownEditor';
import rehypeSanitize from "rehype-sanitize";
import MDEditor, { commands } from "@uiw/react-md-editor";


interface PostCreatorProps {
  isLoggedIn: boolean;
  onCreatePost: (content: string) => Promise<void>;
}

const PostCreator: React.FC<PostCreatorProps> = ({ isLoggedIn, onCreatePost }) => {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isAIModalOpen, setIsAIModalOpen] = useState(false);

  const maxLength = 500;
  const remainingChars = maxLength - content.length;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim() || !isLoggedIn || isSubmitting) return;

    setIsSubmitting(true);
    try {
      await onCreatePost(content.trim());
      setContent('');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    setContent('');
  };

  const handleOpenAITool = () => {
    setIsAIModalOpen(true);
  };

  const handleCloseAITool = () => {
    setIsAIModalOpen(false);
  };

  const handleUseAIContent = (aiContent: string) => {
    setContent(aiContent);
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 mb-6">
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <MDEditor
            value={content}
            autoFocus={true}
            autoFocusEnd={true}
            visibleDragbar={false}
            preview="edit"
            height="100%"
            onChange={value => isLoggedIn && setContent(value || '')}
            previewOptions={{
              rehypePlugins: [[rehypeSanitize]],
            }}
            extraCommands={[commands.codeEdit, commands.codeLive, commands.fullscreen]}
            textareaProps={{
              placeholder: isLoggedIn ? "分享你的想法..." : "請先登入以發布貼文",
              maxLength: maxLength
            }}
          />
          {/* <MarkdownEditor
            value={content}
            onChange={(value) => setContent(value)}
            disabled={!isLoggedIn}
            maxLength={maxLength}
            placeholder={isLoggedIn ? "分享你的想法..." : "請先登入以發布貼文"}
          /> */}
        </div>

        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-2 sm:space-y-0">
          <div className="flex items-center space-x-4">
            <span className={`text-sm ${remainingChars < 0
              ? 'text-red-500'
              : remainingChars < 50
                ? 'text-yellow-500'
                : 'text-gray-500 dark:text-gray-400'
              }`}>
              {remainingChars} 字元剩餘
            </span>

            <button
              type="button"
              onClick={handleOpenAITool}
              disabled={isSubmitting}
              className="flex items-center space-x-1 px-3 py-1 text-purple-600 dark:text-purple-400 hover:text-purple-800 dark:hover:text-purple-300 text-sm font-medium transition-colors disabled:opacity-50 bg-purple-50 dark:bg-purple-900/20 hover:bg-purple-100 dark:hover:bg-purple-900/30 rounded-md"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <span>AI 工具</span>
            </button>
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
      </form>

      <AIToolModal
        isOpen={isAIModalOpen}
        onClose={handleCloseAITool}
        currentContent={content}
        onUseAIContent={handleUseAIContent}
      />
    </div>
  );
};

export default PostCreator;