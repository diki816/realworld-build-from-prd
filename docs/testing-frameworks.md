# Testing Frameworks Documentation

## Vitest (Frontend Testing)

### Version Information
- **Current Stable**: Vitest 3.0.x
- **Performance**: 2-5x faster than traditional testing frameworks
- **Node.js Compatibility**: Works with Node.js 22+
- **TypeScript**: First-class TypeScript support out of the box

### Official Documentation
- **Main Documentation**: https://vitest.dev/
- **Getting Started**: https://vitest.dev/guide/
- **Configuration**: https://vitest.dev/config/
- **API Reference**: https://vitest.dev/api/
- **Browser Mode**: https://vitest.dev/guide/browser/
- **UI Mode**: https://vitest.dev/guide/ui.html
- **GitHub Repository**: https://github.com/vitest-dev/vitest

### Key Features for RealWorld
- **Vite Integration**: Native integration with Vite configuration
- **Jest Compatibility**: 85-90% Jest API compatibility for easy migration
- **ESM Support**: Native ES modules support
- **Browser Testing**: Run tests in real browsers
- **Watch Mode**: Intelligent test re-running based on file changes
- **Snapshot Testing**: Component and data snapshot testing
- **Coverage Reports**: Built-in code coverage with c8

### Installation & Setup
```bash
# Install Vitest and testing utilities
npm install -D vitest@^3.0.0
npm install -D @testing-library/react@^16.0.0
npm install -D @testing-library/jest-dom@^6.0.0
npm install -D @testing-library/user-event@^14.0.0
npm install -D jsdom@^25.0.0
```

### Vitest Configuration
```typescript
// vitest.config.ts
import { defineConfig, mergeConfig } from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      // Test environment
      environment: 'jsdom',
      
      // Setup files
      setupFiles: ['./src/test/setup.ts'],
      
      // Global test utilities
      globals: true,
      
      // Coverage configuration
      coverage: {
        provider: 'v8',
        reporter: ['text', 'html', 'json'],
        exclude: [
          'node_modules/',
          'src/test/',
          '**/*.config.*',
          '**/*.d.ts',
        ],
        thresholds: {
          global: {
            branches: 80,
            functions: 80,
            lines: 80,
            statements: 80
          }
        }
      },
      
      // Test patterns
      include: ['src/**/*.{test,spec}.{ts,tsx}'],
      exclude: ['node_modules/', 'dist/', 'e2e/'],
      
      // Watch mode configuration
      watchExclude: ['node_modules/', 'dist/'],
      
      // Test timeout
      testTimeout: 10000,
      
      // Browser mode (optional)
      // browser: {
      //   enabled: true,
      //   name: 'chromium',
      //   provider: 'playwright'
      // }
    }
  })
)
```

### Test Setup
```typescript
// src/test/setup.ts
import { expect, afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'
import * as matchers from '@testing-library/jest-dom/matchers'

// Extend Vitest's expect with jest-dom matchers
expect.extend(matchers)

// Cleanup after each test
afterEach(() => {
  cleanup()
})

// Mock IntersectionObserver
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  unobserve() {}
}

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  constructor() {}
  disconnect() {}
  observe() {}
  unobserve() {}
}

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
```

### RealWorld Component Testing Examples
```typescript
// src/components/ArticleCard.test.tsx
import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ArticleCard } from './ArticleCard'
import type { Article } from '@/types/article'

const mockArticle: Article = {
  id: 1,
  slug: 'test-article',
  title: 'Test Article',
  description: 'A test article description',
  body: 'Test article body',
  tagList: ['test', 'vitest'],
  createdAt: '2025-01-01T00:00:00.000Z',
  updatedAt: '2025-01-01T00:00:00.000Z',
  favorited: false,
  favoritesCount: 5,
  author: {
    username: 'testuser',
    bio: 'Test user bio',
    image: 'https://example.com/avatar.jpg',
    following: false
  }
}

const createTestQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false }
  }
})

function renderWithProviders(component: React.ReactElement) {
  const queryClient = createTestQueryClient()
  return render(
    <QueryClientProvider client={queryClient}>
      {component}
    </QueryClientProvider>
  )
}

describe('ArticleCard', () => {
  it('renders article information correctly', () => {
    renderWithProviders(<ArticleCard article={mockArticle} />)
    
    expect(screen.getByText('Test Article')).toBeInTheDocument()
    expect(screen.getByText('A test article description')).toBeInTheDocument()
    expect(screen.getByText('testuser')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument() // favorites count
    expect(screen.getByText('test')).toBeInTheDocument()
    expect(screen.getByText('vitest')).toBeInTheDocument()
  })

  it('calls onFavorite when favorite button is clicked', async () => {
    const mockOnFavorite = vi.fn()
    
    renderWithProviders(
      <ArticleCard article={mockArticle} onFavorite={mockOnFavorite} />
    )
    
    const favoriteButton = screen.getByRole('button', { name: /5/ })
    fireEvent.click(favoriteButton)
    
    expect(mockOnFavorite).toHaveBeenCalledWith('test-article')
  })

  it('shows favorited state correctly', () => {
    const favoritedArticle = { ...mockArticle, favorited: true }
    
    renderWithProviders(<ArticleCard article={favoritedArticle} />)
    
    const favoriteButton = screen.getByRole('button', { name: /5/ })
    expect(favoriteButton).toHaveClass('text-red-500')
  })
})

// src/hooks/useArticles.test.ts
import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useArticles } from './useArticles'
import { api } from '@/lib/api'

vi.mock('@/lib/api')
const mockApi = vi.mocked(api)

describe('useArticles', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false }
      }
    })
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )

  it('fetches articles successfully', async () => {
    const mockResponse = {
      articles: [mockArticle],
      articlesCount: 1
    }
    
    mockApi.getArticles.mockResolvedValue(mockResponse)
    
    const { result } = renderHook(() => useArticles(), { wrapper })
    
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })
    
    expect(result.current.data).toEqual(mockResponse)
    expect(mockApi.getArticles).toHaveBeenCalledWith({})
  })

  it('handles filters correctly', async () => {
    const filters = { tag: 'react', limit: 10 }
    
    mockApi.getArticles.mockResolvedValue({
      articles: [],
      articlesCount: 0
    })
    
    const { result } = renderHook(() => useArticles(filters), { wrapper })
    
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })
    
    expect(mockApi.getArticles).toHaveBeenCalledWith(filters)
  })
})
```

### Integration Testing with MSW
```typescript
// src/test/mocks/handlers.ts
import { http, HttpResponse } from 'msw'
import type { Article } from '@/types/article'

const mockArticles: Article[] = [
  // ... mock data
]

export const handlers = [
  // Get articles
  http.get('/api/articles', ({ request }) => {
    const url = new URL(request.url)
    const limit = Number(url.searchParams.get('limit')) || 20
    const offset = Number(url.searchParams.get('offset')) || 0
    
    const paginatedArticles = mockArticles.slice(offset, offset + limit)
    
    return HttpResponse.json({
      articles: paginatedArticles,
      articlesCount: mockArticles.length
    })
  }),

  // Get single article
  http.get('/api/articles/:slug', ({ params }) => {
    const article = mockArticles.find(a => a.slug === params.slug)
    
    if (!article) {
      return HttpResponse.json(
        { errors: { body: ['Article not found'] } },
        { status: 404 }
      )
    }
    
    return HttpResponse.json({ article })
  }),

  // Create article
  http.post('/api/articles', async ({ request }) => {
    const body = await request.json()
    const newArticle: Article = {
      // ... create mock article from request body
    }
    
    return HttpResponse.json({ article: newArticle }, { status: 201 })
  })
]

// src/test/mocks/server.ts
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)

// src/test/setup.ts (addition)
import { server } from './mocks/server'

// Start server before all tests
beforeAll(() => server.listen())

// Reset handlers after each test
afterEach(() => server.resetHandlers())

// Close server after all tests
afterAll(() => server.close())
```

## Playwright (E2E Testing)

### Version Information
- **Current Stable**: Playwright 1.55.0
- **Browser Support**: Chromium, Firefox, WebKit
- **Platform Support**: Windows, macOS, Linux, Docker
- **CI/CD Integration**: Excellent support for GitHub Actions, GitLab CI, Azure DevOps

### Official Documentation
- **Main Documentation**: https://playwright.dev/
- **Getting Started**: https://playwright.dev/docs/intro
- **API Reference**: https://playwright.dev/docs/api/class-playwright
- **Test Generator**: https://playwright.dev/docs/codegen
- **CI/CD Guide**: https://playwright.dev/docs/ci
- **Best Practices**: https://playwright.dev/docs/best-practices
- **GitHub Repository**: https://github.com/microsoft/playwright

### Key Features for RealWorld
- **Cross-Browser Testing**: Test in Chromium, Firefox, and WebKit
- **Auto-Waiting**: Intelligent waiting for elements and network requests
- **Network Interception**: Mock API responses for consistent testing
- **Visual Comparisons**: Screenshot and visual regression testing
- **Mobile Testing**: Device emulation for responsive testing
- **Parallel Execution**: Run tests in parallel for faster feedback

### Installation & Setup
```bash
# Install Playwright
npm install -D @playwright/test@^1.55.0

# Install browsers
npx playwright install

# Install system dependencies (Linux)
npx playwright install-deps
```

### Playwright Configuration
```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  
  // Run tests in files in parallel
  fullyParallel: true,
  
  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,
  
  // Retry on CI only
  retries: process.env.CI ? 2 : 0,
  
  // Opt out of parallel tests on CI
  workers: process.env.CI ? 1 : undefined,
  
  // Reporter configuration
  reporter: process.env.CI 
    ? [['github'], ['html']]
    : [['list'], ['html']],
  
  // Global test configuration
  use: {
    // Base URL for tests
    baseURL: 'http://localhost:3000',
    
    // Collect trace when retrying the failed test
    trace: 'on-first-retry',
    
    // Take screenshot only when test fails
    screenshot: 'only-on-failure',
    
    // Record video only when test fails
    video: 'retain-on-failure',
    
    // Navigation timeout
    navigationTimeout: 30000,
    
    // Action timeout
    actionTimeout: 10000
  },

  // Configure projects for major browsers
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    
    // Mobile browsers
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },
  ],

  // Run your local dev server before starting the tests
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120000
  },
})
```

### RealWorld E2E Test Examples
```typescript
// e2e/auth.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('user can register successfully', async ({ page }) => {
    // Navigate to register page
    await page.click('text=Sign up')
    await expect(page).toHaveURL('/#/register')
    
    // Fill registration form
    const username = `testuser${Date.now()}`
    const email = `${username}@example.com`
    
    await page.fill('input[placeholder="Username"]', username)
    await page.fill('input[placeholder="Email"]', email)
    await page.fill('input[placeholder="Password"]', 'password123')
    
    // Submit form
    await page.click('button[type="submit"]')
    
    // Verify successful registration (redirect to home page)
    await expect(page).toHaveURL('/#/')
    
    // Verify user is logged in
    await expect(page.locator(`text=${username}`)).toBeVisible()
    await expect(page.locator('text=New Article')).toBeVisible()
  })

  test('user can login successfully', async ({ page }) => {
    // Navigate to login page
    await page.click('text=Sign in')
    await expect(page).toHaveURL('/#/login')
    
    // Fill login form with existing user
    await page.fill('input[placeholder="Email"]', 'demo@realworld.io')
    await page.fill('input[placeholder="Password"]', 'password')
    
    // Submit form
    await page.click('button[type="submit"]')
    
    // Verify successful login
    await expect(page).toHaveURL('/#/')
    await expect(page.locator('text=demo')).toBeVisible()
    await expect(page.locator('text=New Article')).toBeVisible()
  })

  test('shows validation errors for invalid login', async ({ page }) => {
    await page.click('text=Sign in')
    
    // Try to submit with empty fields
    await page.click('button[type="submit"]')
    
    // Verify validation errors are shown
    await expect(page.locator('text=Email is required')).toBeVisible()
    await expect(page.locator('text=Password is required')).toBeVisible()
    
    // Try with invalid credentials
    await page.fill('input[placeholder="Email"]', 'invalid@example.com')
    await page.fill('input[placeholder="Password"]', 'wrongpassword')
    await page.click('button[type="submit"]')
    
    await expect(page.locator('text=Invalid credentials')).toBeVisible()
  })
})

// e2e/articles.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Articles', () => {
  // Login before each test
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.click('text=Sign in')
    await page.fill('input[placeholder="Email"]', 'demo@realworld.io')
    await page.fill('input[placeholder="Password"]', 'password')
    await page.click('button[type="submit"]')
    await expect(page).toHaveURL('/#/')
  })

  test('user can create a new article', async ({ page }) => {
    // Navigate to editor
    await page.click('text=New Article')
    await expect(page).toHaveURL('/#/editor')
    
    // Fill article form
    const articleTitle = `Test Article ${Date.now()}`
    await page.fill('input[placeholder="Article Title"]', articleTitle)
    await page.fill('input[placeholder="What\'s this article about?"]', 'Test description')
    await page.fill('textarea[placeholder="Write your article (in markdown)"]', '# Test Content\n\nThis is a test article.')
    await page.fill('input[placeholder="Enter tags"]', 'test')
    await page.keyboard.press('Enter')
    
    // Publish article
    await page.click('button:has-text("Publish Article")')
    
    // Verify article was created
    await expect(page).toHaveURL(/\/#\/article\/.*/)
    await expect(page.locator('h1')).toContainText(articleTitle)
    await expect(page.locator('text=Test Content')).toBeVisible()
    await expect(page.locator('.tag-default:has-text("test")')).toBeVisible()
  })

  test('user can favorite an article', async ({ page }) => {
    // Go to an article
    await page.click('.article-preview:first-child a')
    
    // Get initial favorite count
    const favoriteButton = page.locator('button:has-text("♥")')
    const initialCount = await favoriteButton.textContent()
    
    // Click favorite button
    await favoriteButton.click()
    
    // Verify favorite count increased
    await expect(favoriteButton).not.toContainText(initialCount!)
    
    // Verify button style changed (favorited state)
    await expect(favoriteButton).toHaveClass(/btn-primary/)
  })

  test('user can comment on an article', async ({ page }) => {
    // Go to an article
    await page.click('.article-preview:first-child a')
    
    // Scroll to comments section
    await page.locator('.card:has-text("Write a comment")').scrollIntoViewIfNeeded()
    
    // Write a comment
    const commentText = `Test comment ${Date.now()}`
    await page.fill('textarea[placeholder="Write a comment..."]', commentText)
    await page.click('button:has-text("Post Comment")')
    
    // Verify comment appears
    await expect(page.locator(`.card:has-text("${commentText}")`)).toBeVisible()
  })

  test('article feed pagination works correctly', async ({ page }) => {
    // Verify initial articles are loaded
    await expect(page.locator('.article-preview')).toHaveCount(10)
    
    // Click next page if available
    const nextButton = page.locator('nav[aria-label="Pagination"] button:last-child')
    if (await nextButton.isEnabled()) {
      await nextButton.click()
      
      // Verify URL includes offset parameter
      await expect(page).toHaveURL(/offset=/)
      
      // Verify new articles are loaded
      await expect(page.locator('.article-preview')).toHaveCount(10)
    }
  })
})

// e2e/visual.spec.ts - Visual regression testing
test.describe('Visual Tests', () => {
  test('home page looks correct', async ({ page }) => {
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    
    // Take full page screenshot
    await expect(page).toHaveScreenshot('homepage.png', {
      fullPage: true
    })
  })

  test('article page looks correct', async ({ page }) => {
    await page.goto('/')
    await page.click('.article-preview:first-child a')
    await page.waitForLoadState('networkidle')
    
    // Take screenshot of article content
    await expect(page.locator('.article-content')).toHaveScreenshot('article-content.png')
  })
})
```

## Go Testing (Backend)

### Built-in Testing Package
```go
// internal/handlers/articles_test.go
package handlers

import (
    "bytes"
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    
    // Run migrations
    schema := `
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username VARCHAR(255) UNIQUE NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL
        );
        
        CREATE TABLE articles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            slug VARCHAR(255) UNIQUE NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            body TEXT NOT NULL,
            author_id INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (author_id) REFERENCES users(id)
        );
    `
    
    _, err = db.Exec(schema)
    require.NoError(t, err)
    
    return db
}

func TestCreateArticle(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    handler := &Handler{
        DB:        db,
        JWTSecret: "test-secret",
    }
    
    // Create test user
    _, err := db.Exec(
        "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
        "testuser", "test@example.com", "hashedpassword",
    )
    require.NoError(t, err)
    
    tests := []struct {
        name           string
        requestBody    interface{}
        expectedStatus int
        expectedSlug   string
    }{
        {
            name: "valid article creation",
            requestBody: map[string]interface{}{
                "article": map[string]interface{}{
                    "title":       "Test Article",
                    "description": "Test description",
                    "body":        "Test body",
                    "tagList":     []string{"test"},
                },
            },
            expectedStatus: 201,
            expectedSlug:   "test-article",
        },
        {
            name: "missing title",
            requestBody: map[string]interface{}{
                "article": map[string]interface{}{
                    "description": "Test description",
                    "body":        "Test body",
                },
            },
            expectedStatus: 400,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/api/articles", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            // Add user context (simulate authentication middleware)
            ctx := context.WithValue(req.Context(), "user", &User{ID: 1, Username: "testuser"})
            req = req.WithContext(ctx)
            
            rr := httptest.NewRecorder()
            handler.CreateArticle(rr, req)
            
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedStatus == 201 {
                var response map[string]interface{}
                err := json.Unmarshal(rr.Body.Bytes(), &response)
                require.NoError(t, err)
                
                article := response["article"].(map[string]interface{})
                assert.Equal(t, tt.expectedSlug, article["slug"])
                assert.Equal(t, "Test Article", article["title"])
            }
        })
    }
}

// Benchmark test
func BenchmarkGetArticles(b *testing.B) {
    db := setupTestDB(&testing.T{})
    defer db.Close()
    
    handler := &Handler{DB: db}
    
    // Populate test data
    for i := 0; i < 100; i++ {
        _, err := db.Exec(
            "INSERT INTO articles (slug, title, body, author_id) VALUES (?, ?, ?, ?)",
            fmt.Sprintf("article-%d", i),
            fmt.Sprintf("Article %d", i),
            "Test body",
            1,
        )
        if err != nil {
            b.Fatal(err)
        }
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("GET", "/api/articles", nil)
        rr := httptest.NewRecorder()
        handler.ListArticles(rr, req)
        
        if rr.Code != 200 {
            b.Fatalf("Expected status 200, got %d", rr.Code)
        }
    }
}
```

### Test Coverage and CI Integration
```bash
# Package.json scripts for comprehensive testing
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:run": "vitest run",
    "test:coverage": "vitest run --coverage",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:headed": "playwright test --headed",
    "test:all": "npm run test:run && npm run test:e2e"
  }
}

# Go testing commands
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage
go test -race ./...              # Race condition detection
go test -bench=. ./...           # Run benchmarks
go test -coverprofile=coverage.out ./...  # Generate coverage file
go tool cover -html=coverage.out  # View coverage in browser
```

This comprehensive testing setup ensures high-quality, reliable code for the RealWorld application with frontend unit/integration tests, E2E tests, and backend API testing.