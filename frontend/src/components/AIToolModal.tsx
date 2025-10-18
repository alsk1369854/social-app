import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import AIAPI from '../apis/ai';
import MarkdownEditor from './MarkdownEditor';

interface AIToolModalProps {
  isOpen: boolean;
  onClose: () => void;
  currentContent: string;
  onUseAIContent: (content: string) => void;
}

type AIToolMode = 'create' | 'optimize';

const AIToolModal: React.FC<AIToolModalProps> = ({
  isOpen,
  onClose,
  currentContent,
  onUseAIContent,
}) => {
  const { state } = useAuth();
  const [mode, setMode] = useState<AIToolMode>('create');
  const [topic, setTopic] = useState('');
  const [style, setStyle] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [aiContent, setAiContent] = useState('');

  useEffect(() => {
    if (!isOpen) {
      // 重置狀態
      setTopic('');
      setStyle('');
      setAiContent('');
      setError(null);
      setMode('create');
      setIsLoading(false);
    }
  }, [isOpen]);


  const handleClose = () => {
    onClose();
  };

  const handleUseContent = () => {
    if (aiContent.trim()) {
      onUseAIContent(aiContent.trim());
      onClose();
    }
  };

  const handleGenerate = async () => {
    if (!state.accessToken) {
      setError('請先登入以使用AI工具');
      return;
    }

    setIsLoading(true);
    setError(null);
    setAiContent('');

    const onChunk = (chunk: string) => { setAiContent((prev) => prev + chunk) }
    const onComplete = () => { setIsLoading(false) }
    const onError = (streamError: string) => {
      setIsLoading(false);
      setError(streamError);
    }
    if (mode === 'create') {
      await AIAPI.createPostContentStream(
        { topic: topic.trim(), style: style.trim() },
        state.accessToken,
        onChunk,
        onComplete,
        onError
      );
    } else {
      await AIAPI.createPostContentStream(
        { topic: currentContent.trim(), style: style.trim() },
        state.accessToken,
        onChunk,
        onComplete,
        onError
      );
    }

    setIsLoading(false);
  };


  const canGenerate = () => {
    if (mode === 'create') {
      return topic.trim() && style.trim();
    } else {
      return style.trim() && currentContent.trim();
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
              AI 創作工具
            </h2>
            <button
              onClick={handleClose}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div className="mb-6">
            <div className="flex space-x-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
              <button
                onClick={() => setMode('create')}
                className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-colors ${mode === 'create'
                  ? 'bg-white dark:bg-gray-600 text-blue-600 dark:text-blue-400 shadow-sm'
                  : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
                  }`}
              >
                文章創作
              </button>
              <button
                onClick={() => setMode('optimize')}
                className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-colors ${mode === 'optimize'
                  ? 'bg-white dark:bg-gray-600 text-blue-600 dark:text-blue-400 shadow-sm'
                  : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
                  }`}
              >
                內容轉換
              </button>
            </div>
          </div>

          <div className="space-y-4 mb-6">
            {mode === 'create' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  主題
                </label>
                <input
                  type="text"
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                  placeholder="請輸入想要創作的主題..."
                  className="w-full p-3 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            )}

            {mode === 'optimize' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  目前內容
                </label>
                <div className="p-3 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                  {currentContent || '目前編輯區沒有內容'}
                </div>
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                風格
              </label>
              <select
                value={style}
                onChange={(e) => setStyle(e.target.value)}
                className="w-full p-3 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">請選擇風格...</option>
                <option value="正式">正式</option>
                <option value="輕鬆">輕鬆</option>
                <option value="專業">專業</option>
                <option value="幽默">幽默</option>
                <option value="友善">友善</option>
                <option value="激勵">激勵</option>
                <option value="學術">學術</option>
                <option value="創意">創意</option>
              </select>
            </div>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md">
              <p className="text-red-600 dark:text-red-400 text-sm">{error}</p>
            </div>
          )}

          {(aiContent || isLoading) && (
            <div className="mb-6">
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  AI 生成內容
                </label>
                {isLoading && (
                  <div className="flex items-center space-x-2 text-xs text-blue-600 dark:text-blue-400">
                    <div className="flex space-x-1">
                      <div className="w-1 h-1 bg-blue-600 dark:bg-blue-400 rounded-full animate-bounce"></div>
                      <div className="w-1 h-1 bg-blue-600 dark:bg-blue-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                      <div className="w-1 h-1 bg-blue-600 dark:bg-blue-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                    </div>
                    <span>正在生成中...</span>
                  </div>
                )}
              </div>
              <div className="relative">
                <MarkdownEditor
                  value={aiContent}
                  onChange={setAiContent}
                  placeholder="AI 生成的內容會在這裡顯示..."
                  disabled={isLoading}
                  autoExpand={true}
                  preventAutoCollapse={true}
                  className={`min-h-32 ${isLoading ? 'cursor-not-allowed bg-gray-50 dark:bg-gray-800 border-blue-200 dark:border-blue-800' : ''
                    }`}
                />
              </div>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {isLoading
                  ? 'AI 正在生成內容，請等待完成後再編輯...'
                  : '你可以在這裡編輯AI生成的內容'
                }
              </p>
            </div>
          )}

          <div className="flex justify-end space-x-3">
            <button
              onClick={handleClose}
              disabled={isLoading}
              className="px-4 py-2 text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 font-medium transition-colors disabled:opacity-50"
            >
              取消
            </button>

            {!aiContent && !isLoading ? (
              <button
                onClick={handleGenerate}
                disabled={!canGenerate() || isLoading}
                className="px-4 py-2 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white rounded-md font-medium transition-colors flex items-center"
              >
                {isLoading && (
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                )}
                開始生成
              </button>
            ) : aiContent && !isLoading ? (
              <div className="flex space-x-2">
                <button
                  onClick={handleGenerate}
                  disabled={!canGenerate()}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white rounded-md font-medium transition-colors"
                >
                  重新生成
                </button>
                <button
                  onClick={handleUseContent}
                  disabled={!aiContent.trim()}
                  className="px-4 py-2 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 text-white rounded-md font-medium transition-colors"
                >
                  使用此內容
                </button>
              </div>
            ) : isLoading ? (
              <button
                onClick={handleGenerate}
                disabled={true}
                className="px-4 py-2 bg-blue-500 disabled:bg-gray-400 text-white rounded-md font-medium transition-colors flex items-center"
              >
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                正在生成中...
              </button>
            ) : null}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AIToolModal;