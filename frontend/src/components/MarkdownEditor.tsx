import React, { useState } from 'react';
import MarkdownRenderer from './MarkdownRenderer';

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  maxLength?: number;
  className?: string;
  disabled?: boolean;
  onClick?: () => void;
}

const MarkdownEditor: React.FC<MarkdownEditorProps> = ({
  value,
  onChange,
  placeholder = "輸入內容...",
  maxLength,
  className = "",
  disabled = false,
  onClick
}) => {
  const [activeTab, setActiveTab] = useState<'write' | 'preview'>('write');
  const [isExpanded, setIsExpanded] = useState(false);

  const handleTabChange = (tab: 'write' | 'preview') => {
    setActiveTab(tab);
    if (tab === 'preview' && !isExpanded) {
      setIsExpanded(true);
    }
  };

  const handleTextareaClick = () => {
    setIsExpanded(true);
    if (onClick) onClick();
  };

  const remainingChars = maxLength ? maxLength - value.length : null;

  return (
    <div className={`border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 ${className}`}>
      {/* Tab Headers */}
      {isExpanded && (
        <div className="border-b border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800 rounded-t-md">
          <div className="flex">
            <button
              type="button"
              onClick={() => handleTabChange('write')}
              className={`px-4 py-2 text-sm font-medium rounded-tl-md transition-colors ${
                activeTab === 'write'
                  ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white border-b border-white dark:border-gray-700'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700'
              }`}
            >
              編輯
            </button>
            <button
              type="button"
              onClick={() => handleTabChange('preview')}
              className={`px-4 py-2 text-sm font-medium transition-colors ${
                activeTab === 'preview'
                  ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white border-b border-white dark:border-gray-700'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700'
              }`}
            >
              預覽
            </button>
          </div>
        </div>
      )}

      {/* Content Area */}
      <div className="relative">
        {activeTab === 'write' || !isExpanded ? (
          <>
            <textarea
              value={value}
              onChange={(e) => onChange(e.target.value)}
              onClick={handleTextareaClick}
              disabled={disabled}
              maxLength={maxLength}
              placeholder={placeholder}
              className={`w-full p-3 sm:p-4 border-0 bg-transparent resize-y text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all ${
                disabled ? 'cursor-not-allowed opacity-50' : ''
              } ${isExpanded ? 'h-32 sm:h-36 rounded-b-md' : 'h-12 sm:h-14 rounded-md'}`}
            />
            {isExpanded && remainingChars !== null && (
              <div className="absolute bottom-2 right-2 text-xs text-gray-400 bg-white dark:bg-gray-700 px-2 py-1 rounded">
                {remainingChars >= 0 ? `${remainingChars} 剩餘` : `超出 ${Math.abs(remainingChars)}`}
              </div>
            )}
          </>
        ) : (
          <div className="p-3 sm:p-4 min-h-32 rounded-b-md">
            {value.trim() ? (
              <MarkdownRenderer content={value} showPreview={true} />
            ) : (
              <div className="text-gray-500 dark:text-gray-400 italic text-sm">
                沒有內容可預覽...
              </div>
            )}
          </div>
        )}
      </div>

      {/* Markdown Help Text */}
      {isExpanded && (
        <div className="border-t border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-800 px-3 py-2 rounded-b-md">
          <div className="text-xs text-gray-500 dark:text-gray-400">
            支援 Markdown 語法：**粗體** *斜體* ~~刪除線~~ `代碼` [連結](url) &gt; 引用 - 列表
          </div>
        </div>
      )}
    </div>
  );
};

export default MarkdownEditor;