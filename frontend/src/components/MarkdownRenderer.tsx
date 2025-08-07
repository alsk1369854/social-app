import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';

interface MarkdownRendererProps {
  content: string;
  className?: string;
  showPreview?: boolean;
}

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ 
  content, 
  className = '',
  showPreview = true 
}) => {
  if (!showPreview) {
    return (
      <div className={`whitespace-pre-wrap ${className}`}>
        {content}
      </div>
    );
  }

  return (
    <div className={`max-w-none ${className}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeSanitize]}
        components={{
          // Customize link behavior
          a: ({ node, ...props }) => (
            <a 
              {...props} 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 underline"
            />
          ),
          // Customize code blocks
          code: ({ node, className, children, ...props }: any) => {
            const isInline = !className?.includes('language-');
            if (isInline) {
              return (
                <code 
                  className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded text-sm font-mono text-gray-800 dark:text-gray-200" 
                  {...props}
                >
                  {children}
                </code>
              );
            }
            return (
              <code 
                className="block bg-gray-100 dark:bg-gray-700 p-3 rounded text-sm font-mono text-gray-800 dark:text-gray-200 overflow-x-auto" 
                {...props}
              >
                {children}
              </code>
            );
          },
          // Customize blockquotes
          blockquote: ({ node, ...props }) => (
            <blockquote 
              className="border-l-4 border-gray-300 dark:border-gray-600 pl-4 italic text-gray-600 dark:text-gray-400" 
              {...props}
            />
          ),
          // Customize tables
          table: ({ node, ...props }) => (
            <div className="overflow-x-auto">
              <table 
                className="min-w-full border-collapse border border-gray-300 dark:border-gray-600" 
                {...props}
              />
            </div>
          ),
          th: ({ node, ...props }) => (
            <th 
              className="border border-gray-300 dark:border-gray-600 bg-gray-100 dark:bg-gray-700 px-4 py-2 text-left font-medium" 
              {...props}
            />
          ),
          td: ({ node, ...props }) => (
            <td 
              className="border border-gray-300 dark:border-gray-600 px-4 py-2" 
              {...props}
            />
          ),
          // Customize headings
          h1: ({ node, ...props }) => (
            <h1 className="text-xl sm:text-2xl font-bold mb-4 text-gray-900 dark:text-white" {...props} />
          ),
          h2: ({ node, ...props }) => (
            <h2 className="text-lg sm:text-xl font-bold mb-3 text-gray-900 dark:text-white" {...props} />
          ),
          h3: ({ node, ...props }) => (
            <h3 className="text-base sm:text-lg font-bold mb-2 text-gray-900 dark:text-white" {...props} />
          ),
          // Customize lists
          ul: ({ node, ...props }) => (
            <ul className="list-disc pl-6 mb-4 space-y-1 text-gray-900 dark:text-white" {...props} />
          ),
          ol: ({ node, ...props }) => (
            <ol className="list-decimal pl-6 mb-4 space-y-1 text-gray-900 dark:text-white" {...props} />
          ),
          // Customize paragraphs
          p: ({ node, ...props }) => (
            <p className="mb-3 text-gray-900 dark:text-white leading-relaxed" {...props} />
          ),
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
};

export default MarkdownRenderer;