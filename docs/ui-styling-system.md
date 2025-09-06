# UI & Styling System Documentation

## Tailwind CSS Framework

### Version Information
- **Current Stable**: Tailwind CSS 4.1.x
- **Performance**: 5x faster full builds, 100x faster incremental builds
- **Browser Support**: Safari 16.4+, Chrome 111+, Firefox 128+
- **Fallback Option**: Tailwind CSS 3.4.x for older browser support

### Official Documentation
- **Main Documentation**: https://tailwindcss.com/docs
- **Installation Guide**: https://tailwindcss.com/docs/installation
- **Configuration**: https://tailwindcss.com/docs/configuration
- **Utility Classes**: https://tailwindcss.com/docs/utility-first
- **Customization**: https://tailwindcss.com/docs/theme
- **Responsive Design**: https://tailwindcss.com/docs/responsive-design
- **Dark Mode**: https://tailwindcss.com/docs/dark-mode
- **GitHub Repository**: https://github.com/tailwindlabs/tailwindcss

### Key Features for RealWorld
- **Utility-First**: Compose designs directly in HTML/JSX
- **Responsive Design**: Mobile-first responsive breakpoints
- **Dark Mode**: Built-in dark mode support with class strategy
- **Custom Design System**: Design tokens for colors, spacing, typography
- **Performance**: Optimized CSS output with unused style removal
- **Developer Experience**: IntelliSense support in VS Code

### Installation & Setup
```bash
npm install -D tailwindcss@^4.1.0 postcss@latest autoprefixer@latest
npx tailwindcss init -p
```

### Tailwind Configuration for RealWorld
```typescript
// tailwind.config.ts
import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
  ],
  
  darkMode: ['class'], // Enable dark mode with class strategy
  
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    
    extend: {
      // RealWorld brand colors
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        
        // RealWorld specific colors
        conduit: {
          green: '#5CB85C',
          'light-gray': '#F3F3F3',
          'medium-gray': '#999999',
          'dark-gray': '#373A3C',
        }
      },
      
      // Custom border radius
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      
      // Typography scale
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['Fira Code', 'Consolas', 'monospace'],
      },
      
      // Animation enhancements
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
        'fade-in': {
          from: { opacity: '0', transform: 'translateY(-10px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        }
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'fade-in': 'fade-in 0.2s ease-out',
      },
    },
  },
  
  plugins: [
    require('tailwindcss-animate'),
    require('@tailwindcss/typography'), // For markdown content
    require('@tailwindcss/forms'), // Better form styling
  ],
}

export default config
```

### CSS Variables Setup
```css
/* styles/globals.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --card: 0 0% 100%;
    --card-foreground: 222.2 84% 4.9%;
    --popover: 0 0% 100%;
    --popover-foreground: 222.2 84% 4.9%;
    --primary: 142 76% 36%;
    --primary-foreground: 355 7% 97%;
    --secondary: 210 40% 95%;
    --secondary-foreground: 222.2 84% 4.9%;
    --muted: 210 40% 95%;
    --muted-foreground: 215.4 16.3% 46.9%;
    --accent: 210 40% 95%;
    --accent-foreground: 222.2 84% 4.9%;
    --destructive: 0 84.2% 60.2%;
    --destructive-foreground: 210 40% 98%;
    --border: 214.3 31.8% 91.4%;
    --input: 214.3 31.8% 91.4%;
    --ring: 142 76% 36%;
    --radius: 0.5rem;
  }

  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --card: 222.2 84% 4.9%;
    --card-foreground: 210 40% 98%;
    --popover: 222.2 84% 4.9%;
    --popover-foreground: 210 40% 98%;
    --primary: 142 76% 36%;
    --primary-foreground: 355 7% 97%;
    --secondary: 217.2 32.6% 17.5%;
    --secondary-foreground: 210 40% 98%;
    --muted: 217.2 32.6% 17.5%;
    --muted-foreground: 215 20.2% 65.1%;
    --accent: 217.2 32.6% 17.5%;
    --accent-foreground: 210 40% 98%;
    --destructive: 0 62.8% 30.6%;
    --destructive-foreground: 210 40% 98%;
    --border: 217.2 32.6% 17.5%;
    --input: 217.2 32.6% 17.5%;
    --ring: 142 76% 36%;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  
  body {
    @apply bg-background text-foreground;
  }
  
  /* RealWorld specific base styles */
  h1, h2, h3, h4, h5, h6 {
    @apply font-bold tracking-tight;
  }
  
  .prose {
    @apply max-w-none;
  }
}

@layer components {
  /* Custom component classes */
  .btn {
    @apply inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50;
  }
  
  .btn-primary {
    @apply btn bg-primary text-primary-foreground hover:bg-primary/90;
  }
  
  .btn-secondary {
    @apply btn bg-secondary text-secondary-foreground hover:bg-secondary/80;
  }
  
  .input {
    @apply flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50;
  }
}
```

## Shadcn/UI Component System

### Version Information
- **Current Version**: CLI 3.0 with namespaced registries
- **Status**: Stable with enhanced enterprise features
- **Component Count**: 50+ production-ready components
- **Architecture**: Copy-paste components, not a package dependency

### Official Documentation
- **Main Documentation**: https://ui.shadcn.com/
- **Components**: https://ui.shadcn.com/docs/components
- **Installation**: https://ui.shadcn.com/docs/installation/vite
- **Theming**: https://ui.shadcn.com/docs/theming
- **Dark Mode**: https://ui.shadcn.com/docs/dark-mode/vite
- **CLI Reference**: https://ui.shadcn.com/docs/cli
- **GitHub Repository**: https://github.com/shadcn-ui/ui

### Key Features for RealWorld
- **Accessible by Default**: Built on Radix UI primitives
- **Customizable**: Copy-paste and modify components
- **Consistent Design**: Unified design system across components
- **TypeScript First**: Full TypeScript support
- **Framework Agnostic**: Works with React, Vue, Svelte
- **No Runtime Dependencies**: Components are copied into your codebase

### Installation & Setup
```bash
# Initialize shadcn/ui
npx shadcn@latest init

# Add components for RealWorld
npx shadcn@latest add button
npx shadcn@latest add input
npx shadcn@latest add card
npx shadcn@latest add avatar
npx shadcn@latest add badge
npx shadcn@latest add dropdown-menu
npx shadcn@latest add dialog
npx shadcn@latest add form
npx shadcn@latest add textarea
npx shadcn@latest add tabs
npx shadcn@latest add separator
npx shadcn@latest add skeleton
npx shadcn@latest add toast
npx shadcn@latest add alert
```

### Components Configuration
```json
// components.json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.ts",
    "css": "src/styles/globals.css",
    "baseColor": "zinc",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
```

### RealWorld Component Examples
```typescript
// components/ui/article-card.tsx
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Heart, MessageCircle } from "lucide-react"
import type { Article } from "@/types/article"

interface ArticleCardProps {
  article: Article
  onFavorite?: (slug: string) => void
}

export function ArticleCard({ article, onFavorite }: ArticleCardProps) {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Avatar className="h-8 w-8">
              <AvatarImage src={article.author.image} />
              <AvatarFallback>
                {article.author.username.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div>
              <p className="text-sm font-medium">{article.author.username}</p>
              <p className="text-xs text-muted-foreground">
                {new Date(article.createdAt).toLocaleDateString()}
              </p>
            </div>
          </div>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onFavorite?.(article.slug)}
            className={article.favorited ? "text-red-500" : "text-muted-foreground"}
          >
            <Heart className="h-4 w-4 mr-1" />
            {article.favoritesCount}
          </Button>
        </div>
        
        <CardTitle className="line-clamp-2">{article.title}</CardTitle>
        <CardDescription className="line-clamp-2">
          {article.description}
        </CardDescription>
      </CardHeader>
      
      <CardContent>
        <div className="flex items-center justify-between">
          <div className="flex flex-wrap gap-1">
            {article.tagList.map((tag) => (
              <Badge key={tag} variant="secondary" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>
          
          <Button variant="ghost" size="sm">
            <MessageCircle className="h-4 w-4 mr-1" />
            Read more
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

// components/ui/nav-bar.tsx
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Link } from "@tanstack/react-router"
import { User, Settings, LogOut, PenTool } from "lucide-react"

interface NavBarProps {
  user?: User
  onLogout?: () => void
}

export function NavBar({ user, onLogout }: NavBarProps) {
  return (
    <nav className="border-b bg-background">
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between">
          <Link to="/" className="text-xl font-bold text-conduit-green">
            conduit
          </Link>
          
          <div className="flex items-center space-x-4">
            {user ? (
              <>
                <Button variant="ghost" asChild>
                  <Link to="/editor">
                    <PenTool className="h-4 w-4 mr-2" />
                    New Article
                  </Link>
                </Button>
                
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 rounded-full">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={user.image} />
                        <AvatarFallback>
                          {user.username.slice(0, 2).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                    </Button>
                  </DropdownMenuTrigger>
                  
                  <DropdownMenuContent align="end" className="w-56">
                    <DropdownMenuItem asChild>
                      <Link to="/profile/$username" params={{ username: user.username }}>
                        <User className="h-4 w-4 mr-2" />
                        Profile
                      </Link>
                    </DropdownMenuItem>
                    
                    <DropdownMenuItem asChild>
                      <Link to="/settings">
                        <Settings className="h-4 w-4 mr-2" />
                        Settings
                      </Link>
                    </DropdownMenuItem>
                    
                    <DropdownMenuSeparator />
                    
                    <DropdownMenuItem onClick={onLogout}>
                      <LogOut className="h-4 w-4 mr-2" />
                      Log out
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </>
            ) : (
              <>
                <Button variant="ghost" asChild>
                  <Link to="/login">Sign in</Link>
                </Button>
                <Button asChild>
                  <Link to="/register">Sign up</Link>
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
```

## Lucide React Icons

### Version Information
- **Current Stable**: lucide-react 0.460.0+
- **Icon Count**: 1,500+ icons
- **Bundle Size**: Tree-shakeable, only import icons you use
- **Compatibility**: Works seamlessly with Shadcn/UI

### Official Documentation
- **Main Documentation**: https://lucide.dev/
- **Icon Search**: https://lucide.dev/icons/
- **React Guide**: https://lucide.dev/guide/packages/lucide-react
- **Installation**: https://lucide.dev/guide/installation
- **GitHub Repository**: https://github.com/lucide-icons/lucide

### Key Features
- **Tree Shakeable**: Only bundle icons you actually use
- **Consistent Design**: Uniform 24x24 grid with 2px stroke
- **Accessible**: Built-in accessibility attributes
- **Customizable**: Easy to resize and style with Tailwind classes

### Installation & Usage
```bash
npm install lucide-react
```

```typescript
// Import only the icons you need
import { 
  Heart, 
  MessageCircle, 
  User, 
  Settings, 
  LogOut, 
  PenTool,
  Eye,
  EyeOff,
  Search,
  Menu,
  X
} from "lucide-react"

// Usage in components
export function LoginForm() {
  const [showPassword, setShowPassword] = useState(false)
  
  return (
    <div className="space-y-4">
      <div className="relative">
        <Input type="email" placeholder="Email" />
      </div>
      
      <div className="relative">
        <Input 
          type={showPassword ? "text" : "password"} 
          placeholder="Password" 
        />
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="absolute right-0 top-0 h-full px-3"
          onClick={() => setShowPassword(!showPassword)}
        >
          {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
        </Button>
      </div>
      
      <Button type="submit" className="w-full">
        Sign In
      </Button>
    </div>
  )
}
```

## Dark Mode Implementation

### Theme Toggle Component
```typescript
// components/theme-toggle.tsx
import { Moon, Sun } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useTheme } from "@/hooks/use-theme"

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => setTheme(theme === "light" ? "dark" : "light")}
    >
      <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}
```

This comprehensive UI system provides a modern, accessible, and customizable foundation for the RealWorld application interface.