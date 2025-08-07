import React, { useState, useEffect } from 'react';
import MarkdownRenderer from './MarkdownRenderer';

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  maxLength?: number;
  className?: string;
  disabled?: boolean;
  onClick?: () => void;
  autoExpand?: boolean; // Auto-expand when content is provided
  preventAutoCollapse?: boolean; // Prevent auto-collapse on blur
}

const MarkdownEditor: React.FC<MarkdownEditorProps> = ({
  value,
  onChange,
  placeholder = "輸入內容...",
  maxLength,
  className = "",
  disabled = false,
  onClick,
  autoExpand = false,
  preventAutoCollapse = false
}) => {
  const [activeTab, setActiveTab] = useState<'write' | 'preview'>('write');
  const [isExpanded, setIsExpanded] = useState(false);

  // Auto-expand when content is provided and autoExpand is true
  useEffect(() => {
    if (autoExpand && value.trim() && !isExpanded) {
      setIsExpanded(true);
    }
  }, [value, autoExpand, isExpanded]);

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

  const handleTextareaBlur = () => {
    // Don't auto-collapse if preventAutoCollapse is true
    if (preventAutoCollapse) return;
    
    // Use setTimeout to allow for tab clicks and other interactions to complete
    setTimeout(() => {
      // Only collapse if the textarea is empty and we're still on write tab
      if (!value.trim() && activeTab === 'write') {
        setIsExpanded(false);
      }
    }, 150);
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
              onBlur={handleTextareaBlur}
              disabled={disabled}
              maxLength={maxLength}
              placeholder={placeholder}
              className={`w-full p-3 sm:p-4 border-0 bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all ${
                disabled ? 'cursor-not-allowed opacity-50' : ''
              } ${isExpanded ? 'min-h-32 sm:min-h-36 resize-y rounded-b-md' : 'h-12 sm:h-14 resize-none rounded-md'}`}
              style={{ resize: isExpanded ? 'vertical' : 'none' }}
            />
            {isExpanded && remainingChars !== null && (
              <div className="absolute bottom-2 right-2 text-xs text-gray-400 bg-white dark:bg-gray-700 px-2 py-1 rounded">
                {remainingChars >= 0 ? `${remainingChars} 剩餘` : `超出 ${Math.abs(remainingChars)}`}
              </div>
            )}
          </>
        ) : (
          <div className="p-3 sm:p-4 min-h-32 max-h-96 overflow-y-auto rounded-b-md">
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

      {/* Markdown Help Icon */}
      {isExpanded && (
        <div className="border-t border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-800 px-3 py-2 rounded-b-md">
          <div className="flex items-center justify-end">
            <div className="group relative">
              <svg 
                className="w-4 h-4 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 cursor-help transition-colors" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {/* Tooltip */}
              <div className="absolute bottom-full right-0 mb-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
                <div className="bg-gray-900 dark:bg-gray-700 text-white text-xs rounded py-1 px-2 whitespace-nowrap">
                  支援 Markdown：**粗體** *斜體* `代碼` [連結](url) &gt; 引用
                  <div className="absolute top-full right-2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900 dark:border-t-gray-700"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MarkdownEditor;