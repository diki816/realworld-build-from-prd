# TanStack Ecosystem Documentation

## TanStack Query (React Query)

### Version Information
- **Current Stable**: @tanstack/react-query 5.87.1
- **Recommended**: 5.87.x for production
- **Breaking Changes**: v5 introduces stable Suspense support
- **Support**: Active maintenance with regular updates

### Official Documentation
- **Main Documentation**: https://tanstack.com/query/latest
- **Getting Started**: https://tanstack.com/query/latest/docs/framework/react/quick-start
- **API Reference**: https://tanstack.com/query/latest/docs/framework/react/reference
- **Examples**: https://tanstack.com/query/latest/docs/framework/react/examples/simple
- **Migration Guide**: https://tanstack.com/query/v5/docs/framework/react/guides/migrating-to-v5
- **GitHub Repository**: https://github.com/TanStack/query

### Key Features for RealWorld
- **Server State Management**: Caching, synchronization, background updates
- **Optimistic Updates**: Immediate UI feedback before server confirmation
- **Infinite Queries**: Pagination for article feeds
- **Suspense Support**: React Suspense integration with useSuspenseQuery
- **Offline Support**: Cache persistence and offline-first strategies
- **DevTools**: Excellent debugging experience

### Installation
```bash
npm install @tanstack/react-query@^5.87.1
npm install -D @tanstack/react-query-devtools@^5.87.1
```

### Basic Setup for RealWorld
```typescript
// lib/query-client.ts
import { QueryClient } from '@tanstack/react-query'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Stale time: 5 minutes
      staleTime: 5 * 60 * 1000,
      // Garbage collection time: 10 minutes
      gcTime: 10 * 60 * 1000,
      // Retry failed requests 3 times
      retry: 3,
      // Refetch on window focus
      refetchOnWindowFocus: false,
      // Background refetch interval
      refetchInterval: false,
    },
    mutations: {
      // Show error notifications
      onError: (error) => {
        console.error('Mutation error:', error)
      }
    }
  }
})

// App.tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router />
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}
```

### RealWorld API Integration
```typescript
// hooks/useArticles.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'
import type { Article, ArticleFilters } from '@/types/article'

// Fetch articles with filters
export function useArticles(filters: ArticleFilters = {}) {
  return useQuery({
    queryKey: ['articles', filters],
    queryFn: () => api.getArticles(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
    placeholderData: { articles: [], articlesCount: 0 }
  })
}

// Fetch single article
export function useArticle(slug: string) {
  return useQuery({
    queryKey: ['article', slug],
    queryFn: () => api.getArticle(slug),
    enabled: !!slug,
    staleTime: 10 * 60 * 1000 // 10 minutes for individual articles
  })
}

// Create article mutation
export function useCreateArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: api.createArticle,
    onSuccess: (newArticle) => {
      // Invalidate articles list
      queryClient.invalidateQueries({ queryKey: ['articles'] })
      
      // Optimistically add to cache
      queryClient.setQueryData(['article', newArticle.slug], newArticle)
    },
    onError: (error) => {
      console.error('Failed to create article:', error)
    }
  })
}

// Favorite article mutation with optimistic updates
export function useFavoriteArticle() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ slug, favorited }: { slug: string, favorited: boolean }) =>
      favorited ? api.unfavoriteArticle(slug) : api.favoriteArticle(slug),
    
    onMutate: async ({ slug, favorited }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['article', slug] })
      
      // Snapshot previous value
      const previousArticle = queryClient.getQueryData(['article', slug])
      
      // Optimistically update
      queryClient.setQueryData(['article', slug], (old: Article) => ({
        ...old,
        favorited: !favorited,
        favoritesCount: old.favoritesCount + (favorited ? -1 : 1)
      }))
      
      return { previousArticle }
    },
    
    onError: (err, variables, context) => {
      // Rollback on error
      if (context?.previousArticle) {
        queryClient.setQueryData(['article', variables.slug], context.previousArticle)
      }
    },
    
    onSettled: (data, error, variables) => {
      // Always refetch after mutation
      queryClient.invalidateQueries({ queryKey: ['article', variables.slug] })
    }
  })
}

// Infinite query for article feed
export function useInfiniteArticles(filters: ArticleFilters = {}) {
  return useInfiniteQuery({
    queryKey: ['articles', 'infinite', filters],
    queryFn: ({ pageParam = 0 }) => 
      api.getArticles({ ...filters, offset: pageParam, limit: 20 }),
    getNextPageParam: (lastPage, pages) => {
      const totalLoaded = pages.length * 20
      return totalLoaded < lastPage.articlesCount ? totalLoaded : undefined
    },
    initialPageParam: 0
  })
}
```

## TanStack Router

### Version Information
- **Current Stable**: @tanstack/react-router 1.131.35
- **Status**: Production ready as of 2025
- **Type Safety**: 100% type-safe routing with TypeScript
- **Performance**: Built-in caching and search parameter APIs

### Official Documentation
- **Main Documentation**: https://tanstack.com/router/latest
- **Getting Started**: https://tanstack.com/router/latest/docs/framework/react/quick-start
- **API Reference**: https://tanstack.com/router/latest/docs/framework/react/api
- **Examples**: https://tanstack.com/router/latest/docs/framework/react/examples
- **Migration Guide**: https://tanstack.com/router/latest/docs/framework/react/guide/migrate-from-react-router
- **GitHub Repository**: https://github.com/TanStack/router

### Key Features for RealWorld
- **Type-Safe Routing**: Full TypeScript support for routes and params
- **Hash-Based Routing**: Supports `#/path` routing required by RealWorld spec
- **Built-in Caching**: Automatic route-level caching
- **Search Params**: Type-safe search parameter handling
- **Code Splitting**: Automatic route-based code splitting
- **Nested Routes**: Complex route hierarchies with layout routes

### Installation
```bash
npm install @tanstack/react-router@^1.131.35
npm install -D @tanstack/router-cli@^1.131.35
```

### RealWorld Router Setup
```typescript
// routeTree.gen.ts (auto-generated)
import { createRootRoute, createRoute, createRouter } from '@tanstack/react-router'
import { Layout } from '@/components/layout/Layout'
import { HomePage } from '@/pages/HomePage'
import { LoginPage } from '@/pages/LoginPage'
import { RegisterPage } from '@/pages/RegisterPage'
import { ArticlePage } from '@/pages/ArticlePage'
import { ProfilePage } from '@/pages/ProfilePage'
import { EditorPage } from '@/pages/EditorPage'
import { SettingsPage } from '@/pages/SettingsPage'

// Root route with layout
const rootRoute = createRootRoute({
  component: Layout,
  notFoundComponent: () => <div>404 - Page Not Found</div>
})

// Home route
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage
})

// Authentication routes
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
  beforeLoad: ({ context }) => {
    if (context.auth.user) {
      throw redirect({ to: '/' })
    }
  }
})

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  component: RegisterPage
})

// Article routes
const articleRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/article/$slug',
  component: ArticlePage,
  loader: ({ params }) => queryClient.ensureQueryData({
    queryKey: ['article', params.slug],
    queryFn: () => api.getArticle(params.slug)
  })
})

// Editor routes
const editorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/editor',
  component: EditorPage,
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: '/login' })
    }
  }
})

const editArticleRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/editor/$slug',
  component: EditorPage,
  loader: ({ params }) => queryClient.ensureQueryData({
    queryKey: ['article', params.slug],
    queryFn: () => api.getArticle(params.slug)
  })
})

// Profile routes
const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile/$username',
  component: ProfilePage
})

const profileFavoritesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile/$username/favorites',
  component: ProfilePage
})

// Settings route
const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/settings',
  component: SettingsPage,
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: '/login' })
    }
  }
})

// Create route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  registerRoute,
  articleRoute,
  editorRoute,
  editArticleRoute,
  profileRoute,
  profileFavoritesRoute,
  settingsRoute
])

// Router configuration with hash routing
export const router = createRouter({
  routeTree,
  basepath: '#', // Enable hash-based routing for RealWorld spec
  context: {
    auth: {
      user: null,
      login: () => {},
      logout: () => {}
    },
    queryClient
  },
  defaultPreload: 'intent', // Preload on hover/focus
  defaultPreloadStaleTime: 0
})

// Type registration
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
```

### Router Integration
```typescript
// App.tsx
import { RouterProvider } from '@tanstack/react-router'
import { router } from './routeTree.gen'

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}

// Using router in components
import { Link, useNavigate, useParams, useSearch } from '@tanstack/react-router'

export function Navigation() {
  const navigate = useNavigate()
  
  return (
    <nav>
      <Link to="/" activeProps={{ className: 'active' }}>
        Home
      </Link>
      <Link 
        to="/profile/$username" 
        params={{ username: 'johndoe' }}
        activeProps={{ className: 'active' }}
      >
        Profile
      </Link>
    </nav>
  )
}

// Type-safe parameters and search
export function ArticlePage() {
  const { slug } = useParams({ from: '/article/$slug' }) // Type-safe!
  const { tag, page } = useSearch({ from: '/' }) // Type-safe search params
  
  const { data: article } = useArticle(slug)
  
  return <div>{article?.title}</div>
}
```

### Performance Optimizations
```typescript
// Lazy loading routes
const LazyEditorPage = lazy(() => import('@/pages/EditorPage'))

const editorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/editor',
  component: LazyEditorPage,
  pendingComponent: () => <div>Loading editor...</div>
})

// Prefetching data
const articleRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/article/$slug',
  component: ArticlePage,
  loader: ({ params, context }) => 
    context.queryClient.ensureQueryData({
      queryKey: ['article', params.slug],
      queryFn: () => api.getArticle(params.slug),
      staleTime: 10 * 60 * 1000 // 10 minutes
    })
})
```

## Integration Best Practices

### Query Key Patterns
```typescript
// Consistent query key patterns
const queryKeys = {
  all: ['articles'] as const,
  lists: () => [...queryKeys.all, 'list'] as const,
  list: (filters: ArticleFilters) => [...queryKeys.lists(), filters] as const,
  details: () => [...queryKeys.all, 'detail'] as const,
  detail: (slug: string) => [...queryKeys.details(), slug] as const,
  feed: (type: 'global' | 'personal', filters: ArticleFilters = {}) => 
    ['articles', 'feed', type, filters] as const
}
```

### Error Handling
```typescript
// Global error boundary for queries
function QueryErrorBoundary({ children }: { children: React.ReactNode }) {
  return (
    <ErrorBoundary
      fallback={<div>Something went wrong. Please try again.</div>}
      onError={(error) => console.error('Query error:', error)}
    >
      {children}
    </ErrorBoundary>
  )
}
```

This TanStack ecosystem provides type-safe, performant state management and routing perfectly suited for the RealWorld application's requirements.