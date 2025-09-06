import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter, createHashHistory } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import './index.css'

// Import the generated route tree
import { routeTree } from './routeTree.gen'

// Create a hash-based history instance
const hashHistory = createHashHistory()

// Create a new router instance
const router = createRouter({ 
  routeTree, 
  history: hashHistory,
  defaultPreload: 'intent',
})

// Create a client for TanStack Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      gcTime: 1000 * 60 * 10,   // 10 minutes
      retry: (failureCount, error) => {
        // Don't retry on 4xx errors except 408
        if (error instanceof Error) {
          const status = (error as { status?: number }).status
          if (status !== undefined && status >= 400 && status < 500 && status !== 408) {
            return false
          }
        }
        return failureCount < 3
      },
    },
  },
})

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

// Render the app
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </StrictMode>,
  )
}
