import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'
import type {
  LoginRequest,
  RegisterRequest,
  UpdateUserRequest,
  CreateArticleRequest,
  UpdateArticleRequest,
  CreateCommentRequest,
} from '@/types/api'

// Query keys
export const queryKeys = {
  user: ['user'] as const,
  profile: (username: string) => ['profile', username] as const,
  articles: (params?: Record<string, unknown>) => ['articles', params] as const,
  article: (slug: string) => ['article', slug] as const,
  feed: (params?: Record<string, unknown>) => ['feed', params] as const,
  comments: (slug: string) => ['comments', slug] as const,
  tags: ['tags'] as const,
}

// Auth hooks
export function useLogin() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (credentials: LoginRequest) => apiClient.login(credentials),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.user, data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
    },
  })
}

export function useRegister() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (userData: RegisterRequest) => apiClient.register(userData),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.user, data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
    },
  })
}

export function useLogout() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: () => apiClient.logout(),
    onSuccess: () => {
      queryClient.clear()
    },
  })
}

export function useCurrentUser() {
  return useQuery({
    queryKey: queryKeys.user,
    queryFn: () => apiClient.getCurrentUser(),
    enabled: apiClient.isAuthenticated(),
    retry: false,
  })
}

export function useUpdateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (userData: UpdateUserRequest) => apiClient.updateUser(userData),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.user, data)
    },
  })
}

// Profile hooks
export function useProfile(username: string) {
  return useQuery({
    queryKey: queryKeys.profile(username),
    queryFn: () => apiClient.getProfile(username),
    enabled: !!username,
  })
}

export function useFollowUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (username: string) => apiClient.followUser(username),
    onSuccess: (data, username) => {
      queryClient.setQueryData(queryKeys.profile(username), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
    },
  })
}

export function useUnfollowUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (username: string) => apiClient.unfollowUser(username),
    onSuccess: (data, username) => {
      queryClient.setQueryData(queryKeys.profile(username), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
    },
  })
}

// Article hooks
export function useArticles(params?: {
  tag?: string
  author?: string
  favorited?: string
  limit?: number
  offset?: number
}) {
  return useQuery({
    queryKey: queryKeys.articles(params),
    queryFn: () => apiClient.getArticles(params),
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

export function useFeed(limit?: number, offset?: number) {
  return useQuery({
    queryKey: queryKeys.feed({ limit, offset }),
    queryFn: () => apiClient.getFeed(limit, offset),
    enabled: apiClient.isAuthenticated(),
  })
}

export function useArticle(slug: string) {
  return useQuery({
    queryKey: queryKeys.article(slug),
    queryFn: () => apiClient.getArticle(slug),
    enabled: !!slug,
  })
}

export function useCreateArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (articleData: CreateArticleRequest) => apiClient.createArticle(articleData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
      queryClient.invalidateQueries({ queryKey: queryKeys.feed() })
    },
  })
}

export function useUpdateArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ slug, articleData }: { slug: string; articleData: UpdateArticleRequest }) =>
      apiClient.updateArticle(slug, articleData),
    onSuccess: (data, { slug }) => {
      queryClient.setQueryData(queryKeys.article(slug), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
      queryClient.invalidateQueries({ queryKey: queryKeys.feed() })
    },
  })
}

export function useDeleteArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (slug: string) => apiClient.deleteArticle(slug),
    onSuccess: (_, slug) => {
      queryClient.removeQueries({ queryKey: queryKeys.article(slug) })
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
      queryClient.invalidateQueries({ queryKey: queryKeys.feed() })
    },
  })
}

export function useFavoriteArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (slug: string) => apiClient.favoriteArticle(slug),
    onSuccess: (data, slug) => {
      queryClient.setQueryData(queryKeys.article(slug), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
      queryClient.invalidateQueries({ queryKey: queryKeys.feed() })
    },
  })
}

export function useUnfavoriteArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (slug: string) => apiClient.unfavoriteArticle(slug),
    onSuccess: (data, slug) => {
      queryClient.setQueryData(queryKeys.article(slug), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.articles() })
      queryClient.invalidateQueries({ queryKey: queryKeys.feed() })
    },
  })
}

// Comment hooks
export function useComments(slug: string) {
  return useQuery({
    queryKey: queryKeys.comments(slug),
    queryFn: () => apiClient.getComments(slug),
    enabled: !!slug,
  })
}

export function useCreateComment() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ slug, commentData }: { slug: string; commentData: CreateCommentRequest }) =>
      apiClient.createComment(slug, commentData),
    onSuccess: (_, { slug }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.comments(slug) })
    },
  })
}

export function useDeleteComment() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ slug, id }: { slug: string; id: number }) =>
      apiClient.deleteComment(slug, id),
    onSuccess: (_, { slug }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.comments(slug) })
    },
  })
}

// Tag hooks
export function useTags() {
  return useQuery({
    queryKey: queryKeys.tags,
    queryFn: () => apiClient.getTags(),
    staleTime: 1000 * 60 * 10, // 10 minutes
  })
}