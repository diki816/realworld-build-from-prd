-- Seed data for development and testing
-- Insert demo users with bcrypt hashed passwords

-- Demo user (password: password)
INSERT INTO users (username, email, password_hash, bio, image) VALUES
('demo', 'demo@realworld.io', '$2a$10$N9qo8uLOickgx2ZMRZoMye7I6hJgqA5fS5LfDN1EZ7TQfYC2Yce5y', 'Demo user account for testing', 'https://api.realworld.io/images/demo-avatar.png');

-- Admin user (password: admin123)
INSERT INTO users (username, email, password_hash, bio, image) VALUES
('admin', 'admin@realworld.io', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Administrator account', 'https://api.realworld.io/images/admin-avatar.png');

-- Sample user (password: sample123)
INSERT INTO users (username, email, password_hash, bio, image) VALUES
('johndoe', 'john@example.com', '$2a$10$HTLdM6VnKZ0H4EfpUq/7WODVsVjGjZzFH5lEWJjrUu.N8KKgW7N8K', 'Full-stack developer passionate about React and Go', 'https://api.realworld.io/images/smiley-cyrus.jpeg');

-- Popular tags for development
INSERT INTO tags (name) VALUES
('welcome'),
('demo'),
('introduction'),
('getting-started'),
('realworld'),
('react'),
('golang'),
('typescript'),
('javascript'),
('web-development'),
('full-stack'),
('tutorial'),
('programming'),
('coding'),
('software-engineering');

-- Sample article by demo user
INSERT INTO articles (slug, title, description, body, author_id) VALUES
('welcome-to-realworld', 'Welcome to RealWorld!', 'This is your first article in the RealWorld application.', 
'# Welcome to RealWorld!

This is a sample article to demonstrate the RealWorld application functionality.

## What is RealWorld?

RealWorld is a Medium.com clone built to demonstrate real-world application development patterns. Unlike simple "todo" examples, RealWorld shows how to build a production-ready social blogging platform.

## Features

- User registration and authentication
- Article creation and editing
- Comments system
- User following
- Article favoriting
- Tag-based content discovery

## Technology Stack

This implementation uses:

- **Frontend**: React + TypeScript + Vite + TanStack Query + Tailwind CSS
- **Backend**: Go + SQLite + JWT authentication
- **Testing**: Vitest + Playwright + Go testing framework

Happy coding! 🚀', 1);

-- Get the article ID for tagging
-- Note: In a real migration, you'd handle this differently, but for demo purposes:
INSERT INTO article_tags (article_id, tag_id)
SELECT 1, t.id FROM tags t WHERE t.name IN ('welcome', 'demo', 'introduction', 'realworld');

-- Sample comment on the welcome article
INSERT INTO comments (body, author_id, article_id) VALUES
('Great introduction! Looking forward to exploring the codebase.', 3, 1);

-- Demo user favorites the welcome article
INSERT INTO favorites (user_id, article_id) VALUES (1, 1);

-- Sample follow relationship: johndoe follows demo
INSERT INTO follows (follower_id, following_id) VALUES (3, 1);