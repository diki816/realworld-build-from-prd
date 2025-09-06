import type {
  UserResponse,
  LoginRequest,
  RegisterRequest,
  UpdateUserRequest,
  ArticleResponse,
  ArticlesResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
  Comment,
  CommentsResponse,
  CreateCommentRequest,
  ProfileResponse,
  TagsResponse,
  AuthToken,
} from '@/types/api'

// API Configuration
const API_BASE_URL = 'http://localhost:8080/api'

// Token management
const TOKEN_STORAGE_KEY = 'conduit_token'

export class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL
    this.loadToken()
  }

  // Token management methods
  private loadToken(): void {
    const stored = localStorage.getItem(TOKEN_STORAGE_KEY)
    if (stored) {
      try {
        const authToken: AuthToken = JSON.parse(stored)
        if (authToken.expiresAt > Date.now()) {
          this.token = authToken.token
        } else {
          this.clearToken()
        }
      } catch {
        this.clearToken()
      }
    }
  }

  private saveToken(token: string): void {
    const expiresAt = Date.now() + (7 * 24 * 60 * 60 * 1000) // 7 days
    const authToken: AuthToken = { token, expiresAt }
    localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(authToken))
    this.token = token
  }

  private clearToken(): void {
    localStorage.removeItem(TOKEN_STORAGE_KEY)
    this.token = null
  }

  public getToken(): string | null {
    return this.token
  }

  public isAuthenticated(): boolean {
    return this.token !== null
  }

  // HTTP client methods
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    // Add authorization header if token exists
    if (this.token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${this.token}`,
      }
    }

    try {
      const response = await fetch(url, config)
      
      // Handle empty responses
      if (response.status === 204) {
        return {} as T
      }

      const data = await response.json()

      if (!response.ok) {
        // Handle authentication errors
        if (response.status === 401) {
          this.clearToken()
        }
        
        const error = new Error(data.message || `HTTP ${response.status}`) as Error & {
          status: number
          data: unknown
        }
        error.status = response.status
        error.data = data
        throw error
      }

      return data
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error('Network error')
    }
  }

  // Authentication methods
  async login(credentials: LoginRequest): Promise<UserResponse> {
    const response = await this.request<UserResponse>('/users/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    })
    
    // Extract and save token from user object (assuming it's in the response)
    if ('token' in response.user) {
      this.saveToken((response.user as { token: string }).token)
    }
    
    return response
  }

  async register(userData: RegisterRequest): Promise<UserResponse> {
    const response = await this.request<UserResponse>('/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    })
    
    // Extract and save token from user object
    if ('token' in response.user) {
      this.saveToken((response.user as { token: string }).token)
    }
    
    return response
  }

  async logout(): Promise<void> {
    this.clearToken()
  }

  async getCurrentUser(): Promise<UserResponse> {
    return this.request<UserResponse>('/user')
  }

  async updateUser(userData: UpdateUserRequest): Promise<UserResponse> {
    return this.request<UserResponse>('/user', {
      method: 'PUT',
      body: JSON.stringify(userData),
    })
  }

  // Profile methods
  async getProfile(username: string): Promise<ProfileResponse> {
    return this.request<ProfileResponse>(`/profiles/${username}`)
  }

  async followUser(username: string): Promise<ProfileResponse> {
    return this.request<ProfileResponse>(`/profiles/${username}/follow`, {
      method: 'POST',
    })
  }

  async unfollowUser(username: string): Promise<ProfileResponse> {
    return this.request<ProfileResponse>(`/profiles/${username}/follow`, {
      method: 'DELETE',
    })
  }

  // Article methods
  async getArticles(params: {
    tag?: string
    author?: string
    favorited?: string
    limit?: number
    offset?: number
  } = {}): Promise<ArticlesResponse> {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        searchParams.append(key, value.toString())
      }
    })
    
    const query = searchParams.toString()
    const endpoint = query ? `/articles?${query}` : '/articles'
    return this.request<ArticlesResponse>(endpoint)
  }

  async getFeed(limit?: number, offset?: number): Promise<ArticlesResponse> {
    const params = new URLSearchParams()
    if (limit) params.append('limit', limit.toString())
    if (offset) params.append('offset', offset.toString())
    
    const query = params.toString()
    const endpoint = query ? `/articles/feed?${query}` : '/articles/feed'
    return this.request<ArticlesResponse>(endpoint)
  }

  async getArticle(slug: string): Promise<ArticleResponse> {
    return this.request<ArticleResponse>(`/articles/${slug}`)
  }

  async createArticle(articleData: CreateArticleRequest): Promise<ArticleResponse> {
    return this.request<ArticleResponse>('/articles', {
      method: 'POST',
      body: JSON.stringify(articleData),
    })
  }

  async updateArticle(slug: string, articleData: UpdateArticleRequest): Promise<ArticleResponse> {
    return this.request<ArticleResponse>(`/articles/${slug}`, {
      method: 'PUT',
      body: JSON.stringify(articleData),
    })
  }

  async deleteArticle(slug: string): Promise<void> {
    await this.request(`/articles/${slug}`, {
      method: 'DELETE',
    })
  }

  async favoriteArticle(slug: string): Promise<ArticleResponse> {
    return this.request<ArticleResponse>(`/articles/${slug}/favorite`, {
      method: 'POST',
    })
  }

  async unfavoriteArticle(slug: string): Promise<ArticleResponse> {
    return this.request<ArticleResponse>(`/articles/${slug}/favorite`, {
      method: 'DELETE',
    })
  }

  // Comment methods
  async getComments(slug: string): Promise<CommentsResponse> {
    return this.request<CommentsResponse>(`/articles/${slug}/comments`)
  }

  async createComment(slug: string, commentData: CreateCommentRequest): Promise<Comment> {
    return this.request<Comment>(`/articles/${slug}/comments`, {
      method: 'POST',
      body: JSON.stringify(commentData),
    })
  }

  async deleteComment(slug: string, id: number): Promise<void> {
    await this.request(`/articles/${slug}/comments/${id}`, {
      method: 'DELETE',
    })
  }

  // Tag methods
  async getTags(): Promise<TagsResponse> {
    return this.request<TagsResponse>('/tags')
  }
}

// Create and export a default instance
export const apiClient = new ApiClient()

// Export the class for testing purposes
export default ApiClient