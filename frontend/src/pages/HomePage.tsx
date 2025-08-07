import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import Navbar from '../components/Navbar';
import PostCreator from '../components/PostCreator';
import PostsFeed from '../components/PostsFeed';
import { Post, Comment } from '../models/Post';
import PostAPI from '../apis/post';
import { PostGetPostsByKeywordResponseItem } from '../apis/models/post';
import CommentAPI from '../apis/comment';
import { CommentGetListByPostIDResponseItem } from '../apis/models/comment';

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

// Helper function to convert API response to Post model
const convertAPIPostToPost = (apiPost: PostGetPostsByKeywordResponseItem): Post => ({
  id: apiPost.id,
  content: apiPost.content,
  userID: apiPost.author.id,
  username: apiPost.author.username,
  createdAt: apiPost.createdAt,
  updatedAt: apiPost.updatedAt,
  tags: apiPost.tags.map(tag => ({ id: tag.id, name: tag.name }))
});

const HomePage: React.FC = () => {
  const { state, logout } = useAuth();
  const [posts] = useState<Post[]>(mockPosts); // Keep as fallback for initial load
  const [defaultPosts, setDefaultPosts] = useState<Post[]>([]);
  const [postComments, setPostComments] = useState<Record<string, Comment[]>>(mockComments);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Post[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [loading, setLoading] = useState(false);
  const [defaultPostsLoaded, setDefaultPostsLoaded] = useState(false);

  // Use search results when searching, otherwise use default posts from API
  const displayedPosts = searchQuery.trim() ? searchResults : (defaultPostsLoaded ? defaultPosts : posts);

  const loadDefaultPosts = async () => {
    try {
      setLoading(true);
      const response = await PostAPI.searchPosts({
        keyword: '',
        limit: '20',
        offset: '0'
      });
      
      const convertedPosts = response.data.map(convertAPIPostToPost);
      setDefaultPosts(convertedPosts);
      setDefaultPostsLoaded(true);
    } catch (error) {
      console.error('Failed to load default posts:', error);
      // If API fails, use mock posts as fallback
      setDefaultPosts(mockPosts);
      setDefaultPostsLoaded(true);
    } finally {
      setLoading(false);
    }
  };

  // Load default posts on component mount
  useEffect(() => {
    loadDefaultPosts();
  }, []);

  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    
    if (query.trim()) {
      setIsSearching(true);
      try {
        const response = await PostAPI.searchPosts({
          keyword: query,
          limit: '10',
          offset: '0'
        });
        
        const convertedPosts = response.data.map(convertAPIPostToPost);
        setSearchResults(convertedPosts);
      } catch (error) {
        console.error('Search failed:', error);
        setSearchResults([]);
      } finally {
        setIsSearching(false);
      }
    } else {
      setSearchResults([]);
    }
  };

  const handleCreatePost = async (content: string) => {
    if (!state.user || !state.accessToken) return;

    console.log('Creating post with token:', state.accessToken); // Debug token
    console.log('User state:', state.user);
    console.log('Is authenticated:', state.isAuthenticated);
    console.log('LocalStorage token:', localStorage.getItem('accessToken'));
    console.log('LocalStorage user:', localStorage.getItem('user'));
    setLoading(true);
    try {
      // Call the API to create the post (backend will handle tag parsing)
      await PostAPI.createPost({
        content: content, // Send original content with hashtags
        imageURL: '' // Optional, can be extended later
      }, state.accessToken);

      // After successful post creation, refresh the current view
      if (searchQuery.trim()) {
        // If we're currently searching, re-execute the search to get updated results
        console.log('Refreshing search results for query:', searchQuery);
        await handleSearch(searchQuery);
      } else {
        // If we're showing default posts, reload the default posts list
        console.log('Refreshing default posts list');
        await loadDefaultPosts();
      }
    } catch (error) {
      console.error('Failed to create post:', error);
      // Could show an error message to user here
      throw error; // Re-throw so PostCreator can handle it
    } finally {
      setLoading(false);
    }
  };

  const handleAddComment = async (postId: string, content: string) => {
    if (!state.user || !state.accessToken) return;

    try {
      // Call the API to create the comment
      const response = await CommentAPI.createComment({
        content: content,
        postID: postId
      }, state.accessToken);

      console.log('Comment created successfully:', response);

      // Create local Comment object from API response
      const newComment: Comment = {
        id: response.id,
        content: response.content,
        postID: response.postID,
        userID: response.userID,
        username: state.user.username,
        createdAt: new Date().toISOString()
      };

      // Update local comments state
      setPostComments(prevComments => ({
        ...prevComments,
        [postId]: [...(prevComments[postId] || []), newComment]
      }));
    } catch (error) {
      console.error('Failed to create comment:', error);
      throw error; // Re-throw so Post component can handle it
    }
  };

  const handleLoadComments = async (postId: string) => {
    try {
      console.log('Loading comments for post:', postId);
      
      // Call the API to get comments for this post
      const response = await CommentAPI.getCommentsByPostID(postId);
      
      console.log('Comments loaded successfully:', response);

      // Convert API response to Comment objects
      const convertAPICommentToComment = (apiComment: CommentGetListByPostIDResponseItem): Comment => ({
        id: apiComment.id,
        content: apiComment.content,
        postID: apiComment.postID,
        userID: apiComment.userID,
        username: apiComment.userName,
        createdAt: apiComment.createdAt,
        subComments: apiComment.subComments?.map(convertAPICommentToComment) || []
      });

      const comments = response.map(convertAPICommentToComment);
      
      // Update local comments state
      setPostComments(prevComments => ({
        ...prevComments,
        [postId]: comments
      }));
    } catch (error) {
      console.error('Failed to load comments:', error);
      // Keep existing comments or use empty array if first load fails
      setPostComments(prevComments => ({
        ...prevComments,
        [postId]: prevComments[postId] || []
      }));
    }
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
              {isSearching ? (
                <>搜尋中："{searchQuery}"...</>
              ) : (
                <>搜尋結果："{searchQuery}" ({searchResults.length} 則貼文)</>
              )}
            </p>
            <button
              onClick={async () => {
                setSearchQuery('');
                setSearchResults([]);
                // Always refresh default posts when clearing search to ensure latest data
                console.log('Clearing search and refreshing default posts');
                await loadDefaultPosts();
              }}
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
          posts={displayedPosts}
          postComments={postComments}
          isLoggedIn={!!state.user}
          loading={loading || isSearching}
          onAddComment={handleAddComment}
          onLoadComments={handleLoadComments}
        />
      </main>
    </div>
  );
};

export default HomePage;