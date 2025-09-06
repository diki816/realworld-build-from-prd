// User types
export interface User {
  id: number
  username: string
  email: string
  bio: string
  image: string
  following?: boolean
}

export interface UserResponse {
  user: User
}

export interface LoginRequest {
  user: {
    email: string
    password: string
  }
}

export interface RegisterRequest {
  user: {
    username: string
    email: string
    password: string
  }
}

export interface UpdateUserRequest {
  user: {
    username?: string
    email?: string
    password?: string
    bio?: string
    image?: string
  }
}

// Article types
export interface Article {
  id: number
  slug: string
  title: string
  description: string
  body: string
  createdAt: string
  updatedAt: string
  favorited: boolean
  favoritesCount: number
  author: User
  tagList: string[]
}

export interface ArticleResponse {
  article: Article
}

export interface ArticlesResponse {
  articles: Article[]
  articlesCount: number
}

export interface CreateArticleRequest {
  article: {
    title: string
    description: string
    body: string
    tagList: string[]
  }
}

export interface UpdateArticleRequest {
  article: {
    title?: string
    description?: string
    body?: string
    tagList?: string[]
  }
}

// Comment types
export interface Comment {
  id: number
  body: string
  createdAt: string
  updatedAt: string
  author: User
}

export interface CommentResponse {
  comment: Comment
}

export interface CommentsResponse {
  comments: Comment[]
}

export interface CreateCommentRequest {
  comment: {
    body: string
  }
}

// Profile types
export interface ProfileResponse {
  profile: User
}

// Tags types
export interface TagsResponse {
  tags: string[]
}

// Error types
export interface ApiError {
  errors: {
    [key: string]: string[]
  }
}

// Auth types
export interface AuthToken {
  token: string
  expiresAt: number
}