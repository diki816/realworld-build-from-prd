#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script to create GitHub issues based on tasks.md
Run this script after authenticating with GitHub CLI: gh auth login
"""

import subprocess
import json
import re
import sys
import os

# Set UTF-8 encoding for Windows
if sys.platform.startswith('win'):
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

# Define labels that need to be created
REQUIRED_LABELS = [
    {"name": "backend", "description": "Backend development tasks", "color": "0052cc"},
    {"name": "frontend", "description": "Frontend development tasks", "color": "5319e7"},
    {"name": "foundation", "description": "Core foundation tasks", "color": "d93f0b"},
    {"name": "authentication", "description": "Authentication and authorization", "color": "fbca04"},
    {"name": "articles", "description": "Article management features", "color": "0e8a16"},
    {"name": "comments", "description": "Comments system", "color": "006b75"},
    {"name": "tags", "description": "Tags and tagging system", "color": "7057ff"},
    {"name": "setup", "description": "Project setup and configuration", "color": "e99695"},
    {"name": "navigation", "description": "Navigation and routing", "color": "bfd4f2"},
    {"name": "feeds", "description": "Article feeds and listing", "color": "c2e0c6"},
    {"name": "editor", "description": "Article editor functionality", "color": "f9d0c4"},
    {"name": "profiles", "description": "User profiles and settings", "color": "fef2c0"},
    {"name": "ui", "description": "User interface and design", "color": "d4c5f9"},
    {"name": "responsive", "description": "Responsive design and mobile", "color": "c5def5"},
    {"name": "devops", "description": "DevOps and deployment", "color": "0052cc"},
    {"name": "docker", "description": "Docker containerization", "color": "006b75"},
    {"name": "testing", "description": "Testing and quality assurance", "color": "d93f0b"},
    {"name": "quality", "description": "Code quality and standards", "color": "fbca04"},
    {"name": "security", "description": "Security and validation", "color": "b60205"},
    {"name": "error-handling", "description": "Error handling and user feedback", "color": "d93f0b"},
    {"name": "performance", "description": "Performance optimization", "color": "0e8a16"},
    {"name": "documentation", "description": "Documentation and guides", "color": "5319e7"},
    {"name": "deployment", "description": "Production deployment", "color": "0052cc"},
    {"name": "high-priority", "description": "High priority tasks", "color": "b60205"},
    {"name": "medium-priority", "description": "Medium priority tasks", "color": "fbca04"},
]

# Define the task groups based on phases from tasks.md
TASK_GROUPS = [
    {
        "title": "Phase 1.1: Project Setup & Go Backend Foundation",
        "body": """## Tasks:
1. Set up project structure and development environment
2. Initialize Go backend with basic server and SQLite database  
3. Create database schema and migrations for all tables
4. Implement JWT authentication middleware and user models

## Acceptance Criteria:
- [x] Root project structure created as per PRD section 13.1
- [x] Git repository initialized with proper .gitignore
- [x] Go module set up with HTTP server and SQLite connection
- [x] Database schema implemented with all required tables
- [x] JWT authentication middleware working
- [x] User models with validation created

## Priority: High
## Phase: Backend Foundation""",
        "labels": ["backend", "foundation", "high-priority"]
    },
    {
        "title": "Phase 1.2: User Authentication & Profile System",
        "body": """## Tasks:
5. Build user authentication endpoints (register, login, get/update user)
6. Implement user profile endpoints (get profile, follow/unfollow)

## Endpoints to implement:
- `POST /api/users/login` - User login
- `POST /api/users` - User registration
- `GET /api/user` - Get current user (auth required)
- `PUT /api/user` - Update user (auth required)
- `GET /api/profiles/:username` - Get user profile
- `POST /api/profiles/:username/follow` - Follow user (auth required)
- `DELETE /api/profiles/:username/follow` - Unfollow user (auth required)

## Acceptance Criteria:
- [x] All authentication endpoints working with proper validation
- [x] Password hashing with bcrypt implemented
- [x] JWT token handling in responses
- [x] Follow/unfollow functionality working
- [x] Proper error handling and validation

## Priority: High
## Phase: Backend Foundation""",
        "labels": ["backend", "authentication", "high-priority"]
    },
    {
        "title": "Phase 1.3: Articles System - CRUD & Feed",
        "body": """## Tasks:
7. Create article CRUD endpoints with slug generation
8. Build article feed endpoints (global feed, personal feed)
9. Implement article favoriting system

## Endpoints to implement:
- `GET /api/articles/:slug` - Get single article
- `POST /api/articles` - Create article (auth required)
- `PUT /api/articles/:slug` - Update article (auth required)  
- `DELETE /api/articles/:slug` - Delete article (auth required)
- `GET /api/articles` - List articles (with filtering & pagination)
- `GET /api/articles/feed` - Get user feed (auth required)
- `POST /api/articles/:slug/favorite` - Favorite article (auth required)
- `DELETE /api/articles/:slug/favorite` - Unfavorite article (auth required)

## Acceptance Criteria:
- [x] Article CRUD operations working
- [x] Slug generation utility implemented
- [x] Article feeds with filtering by tag, author, favorited status
- [x] Pagination support
- [x] Favoriting system with accurate counts

## Priority: High  
## Phase: Backend Foundation""",
        "labels": ["backend", "articles", "high-priority"]
    },
    {
        "title": "Phase 1.4: Comments & Tags System",
        "body": """## Tasks:
10. Create comments system for articles
11. Build tags system and popular tags endpoint
12. Add CORS middleware and API error handling

## Endpoints to implement:
- `GET /api/articles/:slug/comments` - Get article comments
- `POST /api/articles/:slug/comments` - Add comment (auth required)
- `DELETE /api/articles/:slug/comments/:id` - Delete comment (auth required)
- `GET /api/tags` - Get all tags

## Acceptance Criteria:
- [x] Comments CRUD functionality working
- [x] Tag creation during article creation
- [x] Popular tags functionality
- [x] CORS middleware configured for frontend
- [x] Consistent error response format
- [x] Request logging middleware
- [x] Input validation and sanitization

## Priority: High
## Phase: Backend Foundation""",
        "labels": ["backend", "comments", "tags", "high-priority"]
    },
    {
        "title": "Phase 2.1: Frontend Setup & Architecture",
        "body": """## Tasks:
13. Initialize React frontend with Vite, TypeScript, and Tailwind CSS
14. Set up Shadcn/UI components and design system
15. Configure TanStack Router for hash-based routing
16. Set up TanStack Query for API state management  
17. Create API client with authentication handling

## Acceptance Criteria:
- [x] Vite + React + TypeScript project created
- [x] Tailwind CSS configured with custom settings
- [x] Shadcn/UI initialized with base components
- [x] TanStack Router with hash-based routing configured
- [x] All route definitions created (#/, #/login, #/register, etc.)
- [x] TanStack Query set up with error handling
- [x] API client with JWT token management
- [x] Request/response interceptors working

## Priority: High
## Phase: Frontend Foundation""",
        "labels": ["frontend", "setup", "high-priority"]
    },
    {
        "title": "Phase 2.2: Authentication Pages & Navigation",
        "body": """## Tasks:
18. Build authentication pages (login, register)
19. Implement header navigation with conditional rendering

## Pages to create:
- Login page with email/password validation
- Registration page with username/email/password  
- Header component with logo and navigation
- User dropdown menu for authenticated users

## Acceptance Criteria:
- [x] Login and registration forms with validation
- [x] Form error handling and display
- [x] Navigation between login/register pages
- [x] Header with conditional rendering based on auth status
- [x] User dropdown with logout functionality
- [x] Proper form validation and user feedback

## Priority: High
## Phase: Frontend Foundation""",
        "labels": ["frontend", "authentication", "navigation", "high-priority"]
    },
    {
        "title": "Phase 2.3: Home Page & Article Feeds",
        "body": """## Tasks:
20. Create home page with article feeds and popular tags
21. Build article detail page with comments

## Features to implement:
- Article feed tabs (Your Feed, Global Feed, Tag Feed)
- Article preview cards with metadata
- Popular tags sidebar
- Article detail view with markdown rendering
- Comments section with add/delete functionality
- Follow/favorite buttons
- Author-only edit/delete controls

## Acceptance Criteria:
- [x] Home page with working feed tabs
- [x] Article preview cards showing metadata
- [x] Popular tags sidebar with filtering
- [x] Article detail page with markdown rendering
- [x] Comments system working
- [x] Follow and favorite buttons functional
- [x] Proper author permissions for edit/delete

## Priority: High
## Phase: Frontend Features""",
        "labels": ["frontend", "articles", "feeds", "high-priority"]
    },
    {
        "title": "Phase 2.4: Article Editor & User Profiles",
        "body": """## Tasks:
22. Implement article editor with markdown support
23. Create user profile pages with article tabs
24. Build user settings page for profile updates

## Features to implement:
- Article creation form (title, description, body, tags)
- Markdown editor or textarea
- Tag input with autocomplete
- User profile display with bio and image
- Article tabs: "My Articles" and "Favorited Articles"
- Settings form for profile updates
- Follow button for other users

## Acceptance Criteria:
- [x] Article editor with all required fields
- [x] Tag input with autocomplete functionality
- [x] Publish/update article functionality
- [x] Edit mode for existing articles
- [x] User profile pages with bio and image
- [x] Article tabs working correctly
- [x] Settings page with form validation
- [x] Follow functionality on profiles

## Priority: Medium
## Phase: Frontend Features""",
        "labels": ["frontend", "editor", "profiles", "medium-priority"]
    },
    {
        "title": "Phase 2.5: Pagination & UI Polish",
        "body": """## Tasks:
25. Add pagination components and logic
26. Implement responsive design and mobile optimization

## Features to implement:
- Reusable pagination component
- Page navigation for article feeds
- Loading states for page transitions
- Mobile-first responsive design
- Tablet and desktop breakpoints
- Touch-friendly interface elements

## Acceptance Criteria:
- [x] Pagination working on all feed pages
- [x] Loading states during pagination
- [x] Responsive design across all breakpoints
- [x] Mobile-optimized touch interfaces
- [x] Proper loading states and transitions
- [x] UI polish and smooth interactions

## Priority: Medium
## Phase: Frontend Features""",
        "labels": ["frontend", "ui", "responsive", "medium-priority"]
    },
    {
        "title": "Phase 3.1: Development Environment & Docker",
        "body": """## Tasks:
27. Create Docker configuration for development and production
28. Set up testing infrastructure (backend and frontend)

## Deliverables:
- Dockerfile.dev for backend development
- Dockerfile.dev for frontend development
- docker-compose.yml for full environment
- Production Docker configurations
- Go testing framework setup
- Vitest for frontend unit tests
- React Testing Library for component tests
- Test databases and fixtures

## Acceptance Criteria:
- [x] Docker development environment working
- [x] Docker production builds optimized
- [x] Testing frameworks configured
- [x] Test databases and fixtures set up
- [x] All containers working together
- [x] Easy development workflow

## Priority: Medium
## Phase: Polish & Production""",
        "labels": ["devops", "docker", "testing", "medium-priority"]
    },
    {
        "title": "Phase 3.2: Comprehensive Testing",
        "body": """## Tasks:
29. Write comprehensive tests for API endpoints
30. Add frontend component and integration tests

## Testing scope:
- Unit tests for all handler functions
- Integration tests for database operations  
- API endpoint tests with httptest
- Authentication and authorization tests
- Component tests for all major components
- Integration tests for user flows
- API integration tests with MSW
- E2E tests with Playwright for critical paths

## Acceptance Criteria:
- [x] >80% code coverage for backend
- [x] All API endpoints tested
- [x] Database operations tested
- [x] Frontend components tested
- [x] User flows tested end-to-end
- [x] CI/CD pipeline passing all tests

## Priority: Medium
## Phase: Polish & Production""",
        "labels": ["testing", "quality", "medium-priority"]
    },
    {
        "title": "Phase 3.3: Error Handling & Security",
        "body": """## Tasks:
31. Implement error handling and user feedback
32. Add input validation and security measures

## Features:
- Toast notifications for user actions
- Error boundaries for React
- Network error handling
- Loading states and skeleton screens
- Comprehensive input validation
- XSS protection for user content
- CSRF protection
- Rate limiting for API endpoints

## Acceptance Criteria:
- [x] User-friendly error messages
- [x] Proper error boundaries
- [x] Graceful network error handling
- [x] Security measures implemented
- [x] Input validation on all forms
- [x] Rate limiting configured
- [x] XSS and CSRF protection active

## Priority: High
## Phase: Polish & Production""",
        "labels": ["security", "error-handling", "high-priority"]
    },
    {
        "title": "Phase 3.4: Performance & Production Ready",
        "body": """## Tasks:
33. Optimize performance and add caching strategies
34. Create production deployment configuration
35. Write documentation and setup instructions

## Optimizations:
- Code splitting for frontend
- Image optimization and lazy loading
- API response caching
- Database query optimization
- Production environment variables
- Health check endpoints
- Monitoring and logging setup
- Comprehensive README
- API documentation
- Deployment guides

## Acceptance Criteria:
- [x] Performance optimizations implemented
- [x] Production deployment ready
- [x] Monitoring and logging configured
- [x] Complete documentation written
- [x] Deployment guides created
- [x] API documentation complete

## Priority: Medium
## Phase: Polish & Production""",
        "labels": ["performance", "documentation", "deployment", "medium-priority"]
    }
]

def find_gh_executable():
    """Find GitHub CLI executable in common locations"""
    common_paths = [
        "gh",  # In PATH
        "/c/Program Files/GitHub CLI/gh.exe",  # Common Windows install
        "C:\\Program Files\\GitHub CLI\\gh.exe",  # Windows format
        "/usr/local/bin/gh",  # macOS Homebrew
        "/opt/homebrew/bin/gh",  # macOS Apple Silicon
    ]
    
    for path in common_paths:
        try:
            result = subprocess.run([path, "--version"], capture_output=True, text=True, 
                                  timeout=5, encoding='utf-8', errors='replace')
            if result.returncode == 0:
                return path
        except:
            continue
    return None

def create_github_labels(gh_path):
    """Create GitHub labels for the issues"""
    print("\nCreating GitHub labels...")
    
    for label in REQUIRED_LABELS:
        print(f"Creating label '{label['name']}'...")
        
        cmd = [
            gh_path,
            "label", "create",
            label['name'],
            "--description", label['description'],
            "--color", label['color'],
            "--force"  # This will update existing labels
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=15,
                                  encoding='utf-8', errors='replace')
            
            if result.returncode == 0:
                print(f"[SUCCESS] Label '{label['name']}' created/updated")
            else:
                # Check if it's just because label already exists
                if "already exists" in result.stderr.lower():
                    print(f"[INFO] Label '{label['name']}' already exists")
                else:
                    print(f"[WARNING] Failed to create label '{label['name']}': {result.stderr.strip()}")
                    
        except subprocess.TimeoutExpired:
            print(f"[ERROR] Timeout creating label: {label['name']}")
        except Exception as e:
            print(f"[ERROR] Error creating label '{label['name']}': {e}")

def create_github_issues():
    """Create GitHub issues for each task group"""
    
    # Find GitHub CLI executable
    gh_path = find_gh_executable()
    if not gh_path:
        print("[ERROR] GitHub CLI not found. Please install it from: https://cli.github.com/")
        print("Or manually create issues using the task groups below:")
        print_manual_issues()
        return
    
    print(f"Found GitHub CLI at: {gh_path}")
    
    # First, create all required labels
    create_github_labels(gh_path)
    
    print("\nCreating GitHub issues for RealWorld implementation tasks...")
    
    for i, task_group in enumerate(TASK_GROUPS, 1):
        print(f"\n{i}/{len(TASK_GROUPS)}: Creating issue '{task_group['title']}'")
        
        # Prepare the command
        cmd = [
            gh_path,
            "issue", "create",
            "--title", task_group['title'],
            "--body", task_group['body']
        ]
        
        # Add labels
        for label in task_group['labels']:
            cmd.extend(["--label", label])
        
        try:
            # Execute the command
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=30,
                                  encoding='utf-8', errors='replace')
            
            if result.returncode == 0:
                print(f"[SUCCESS] Issue created successfully: {result.stdout.strip()}")
            else:
                print(f"[ERROR] Failed to create issue: {result.stderr.strip()}")
                
        except subprocess.TimeoutExpired:
            print(f"[ERROR] Timeout creating issue: {task_group['title']}")
        except Exception as e:
            print(f"[ERROR] Error creating issue: {e}")

def print_manual_issues():
    """Print issues for manual creation"""
    print("\n" + "="*80)
    print("MANUAL ISSUE CREATION")
    print("="*80)
    print("First, create these labels in your GitHub repository:\n")
    
    for label in REQUIRED_LABELS:
        print(f"- {label['name']}: {label['description']} (color: #{label['color']})")
    
    print("\nThen copy and paste these issues into GitHub Issues manually:\n")
    
    for i, task_group in enumerate(TASK_GROUPS, 1):
        print(f"--- Issue #{i}: {task_group['title']} ---")
        print(f"Labels: {', '.join(task_group['labels'])}")
        print(f"Body:\n{task_group['body']}")
        print("\n" + "-"*80 + "\n")

if __name__ == "__main__":
    print("GitHub Issues Creation Script")
    print("=" * 50)
    print("This script will create GitHub issues based on the RealWorld implementation tasks.")
    print("Make sure you're authenticated with GitHub CLI first: gh auth login")
    print()
    
    # Check if gh is available and authenticated
    gh_path = find_gh_executable()
    if gh_path:
        try:
            result = subprocess.run([gh_path, "auth", "status"], 
                                  capture_output=True, text=True, timeout=10,
                                  encoding='utf-8', errors='replace')
            if result.returncode != 0:
                print("[ERROR] GitHub CLI not authenticated. Please run: gh auth login")
                print("Authentication is required to create issues.")
                exit(1)
            else:
                print("[SUCCESS] GitHub CLI authenticated")
        except Exception as e:
            print(f"[ERROR] Error checking GitHub CLI auth status: {e}")
            print("Please ensure GitHub CLI is installed and authenticated")
            exit(1)
    else:
        print("[ERROR] GitHub CLI not found.")
        print("Install from: https://cli.github.com/")
        print("Or create issues manually using the output below.")
        print_manual_issues()
        exit(1)
    
    # Confirm before creating issues
    response = input("\nProceed with creating issues? (y/N): ")
    if response.lower() in ['y', 'yes']:
        create_github_issues()
        print("\nDone! Check your repository for the created issues.")
    else:
        print("Cancelled.")