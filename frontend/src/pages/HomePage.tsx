import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import Navbar from '../components/Navbar';
import PostCreator from '../components/PostCreator';
import PostsFeed from '../components/PostsFeed';
import { Post, Comment } from '../models/Post';

// Mock data for demonstration
const mockPosts: Post[] = [
  {
    id: '1',
    content: '歡迎來到 Ming 的社群網站！這是第一則貼文，讓我們開始分享想法吧！',
    userID: '1',
    username: 'DemoUser',
    createdAt: new Date(Date.now() - 3600000).toISOString(), // 1 hour ago
    tags: [{ id: '1', name: '歡迎' }, { id: '2', name: '分享' }]
  },
  {
    id: '2',
    content: '今天天氣真好！大家有什麼有趣的計畫嗎？\n\n期待看到大家的分享 😊',
    userID: '2',
    username: 'TestUser',
    createdAt: new Date(Date.now() - 7200000).toISOString(), // 2 hours ago
    tags: [{ id: '3', name: '天氣' }]
  }
];

const mockComments: Record<string, Comment[]> = {
  '1': [
    {
      id: '1',
      content: '太棒了！期待在這裡和大家交流！',
      postID: '1',
      userID: '2',
      username: 'TestUser',
      createdAt: new Date(Date.now() - 1800000).toISOString() // 30 minutes ago
    },
    {
      id: '2',
      content: '歡迎！希望這個平台能夠蓬勃發展！',
      postID: '1',
      userID: '3',
      username: 'AnotherUser',
      createdAt: new Date(Date.now() - 900000).toISOString() // 15 minutes ago
    }
  ],
  '2': [
    {
      id: '3',
      content: '確實！我計畫去公園散步',
      postID: '2',
      userID: '1',
      username: 'DemoUser',
      createdAt: new Date(Date.now() - 3600000).toISOString() // 1 hour ago
    }
  ]
};

const HomePage: React.FC = () => {
  const { state, logout } = useAuth();
  const [posts, setPosts] = useState<Post[]>(mockPosts);
  const [postComments, setPostComments] = useState<Record<string, Comment[]>>(mockComments);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredPosts, setFilteredPosts] = useState<Post[]>(posts);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    // Filter posts based on search query
    if (searchQuery.trim()) {
      const filtered = posts.filter(post => 
        post.content.toLowerCase().includes(searchQuery.toLowerCase()) ||
        post.username?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        post.tags?.some(tag => tag.name.toLowerCase().includes(searchQuery.toLowerCase()))
      );
      setFilteredPosts(filtered);
    } else {
      setFilteredPosts(posts);
    }
  }, [searchQuery, posts]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const handleCreatePost = async (content: string) => {
    if (!state.user) return;

    setLoading(true);
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    const newPost: Post = {
      id: Date.now().toString(),
      content,
      userID: state.user.id,
      username: state.user.username,
      createdAt: new Date().toISOString(),
      tags: []
    };

    setPosts(prevPosts => [newPost, ...prevPosts]);
    setLoading(false);
  };

  const handleAddComment = async (postId: string, content: string) => {
    if (!state.user) return;

    setLoading(true);
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 500));

    const newComment: Comment = {
      id: Date.now().toString(),
      content,
      postID: postId,
      userID: state.user.id,
      username: state.user.username,
      createdAt: new Date().toISOString()
    };

    setPostComments(prevComments => ({
      ...prevComments,
      [postId]: [...(prevComments[postId] || []), newComment]
    }));
    setLoading(false);
  };

  const handleLoadComments = async (postId: string) => {
    // In a real app, this would fetch comments from the API if they haven't been loaded yet
    // For demo purposes, comments are already loaded
    await new Promise(resolve => setTimeout(resolve, 300));
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
      <Navbar
        user={state.user}
        onSearch={handleSearch}
        onLogoutClick={logout}
      />
      
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        {searchQuery && (
          <div className="mb-4 sm:mb-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 sm:p-4">
            <p className="text-blue-800 dark:text-blue-200 text-sm sm:text-base">
              搜尋結果："{searchQuery}" ({filteredPosts.length} 則貼文)
            </p>
            <button
              onClick={() => setSearchQuery('')}
              className="mt-2 text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200 underline text-sm"
            >
              清除搜尋
            </button>
          </div>
        )}
        
        <PostCreator
          isLoggedIn={!!state.user}
          onCreatePost={handleCreatePost}
        />
        
        <PostsFeed
          posts={filteredPosts}
          postComments={postComments}
          isLoggedIn={!!state.user}
          loading={loading}
          onAddComment={handleAddComment}
          onLoadComments={handleLoadComments}
        />
      </main>
    </div>
  );
};

export default HomePage;