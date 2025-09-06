@echo off
echo Creating GitHub issues for RealWorld implementation...
echo.

REM Check if GitHub CLI is available
gh --version >nul 2>&1
if %errorlevel% neq 0 (
    echo GitHub CLI not found. Please install it first.
    exit /b 1
)

REM Check authentication
gh auth status >nul 2>&1
if %errorlevel% neq 0 (
    echo Please authenticate with GitHub CLI first: gh auth login
    exit /b 1
)

echo Creating Phase 1.1: Project Setup & Go Backend Foundation
gh issue create --title "Phase 1.1: Project Setup & Go Backend Foundation" --body "## Tasks:^

1. Set up project structure and development environment^
2. Initialize Go backend with basic server and SQLite database^
3. Create database schema and migrations for all tables^
4. Implement JWT authentication middleware and user models^

## Acceptance Criteria:^
- ✅ Root project structure created as per PRD section 13.1^
- ✅ Git repository initialized with proper .gitignore^
- ✅ Go module set up with HTTP server and SQLite connection^
- ✅ Database schema implemented with all required tables^
- ✅ JWT authentication middleware working^
- ✅ User models with validation created^

## Priority: High^
## Phase: Backend Foundation" --label "backend,foundation,high-priority"

echo Creating Phase 1.2: User Authentication & Profile System
gh issue create --title "Phase 1.2: User Authentication & Profile System" --body "## Tasks:^

5. Build user authentication endpoints (register, login, get/update user)^
6. Implement user profile endpoints (get profile, follow/unfollow)^

## Endpoints to implement:^
- POST /api/users/login - User login^
- POST /api/users - User registration^
- GET /api/user - Get current user (auth required)^
- PUT /api/user - Update user (auth required)^
- GET /api/profiles/:username - Get user profile^
- POST /api/profiles/:username/follow - Follow user (auth required)^
- DELETE /api/profiles/:username/follow - Unfollow user (auth required)^

## Acceptance Criteria:^
- ✅ All authentication endpoints working with proper validation^
- ✅ Password hashing with bcrypt implemented^
- ✅ JWT token handling in responses^
- ✅ Follow/unfollow functionality working^
- ✅ Proper error handling and validation^

## Priority: High^
## Phase: Backend Foundation" --label "backend,authentication,high-priority"

echo.
echo Issues created successfully! Check your GitHub repository.
pause