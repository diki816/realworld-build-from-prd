# Build Tools Documentation

## Vite Build Tool

### Version Information
- **Current Stable**: Vite 7.1.4
- **Recommended**: 7.1.x for production
- **Node.js Requirement**: 20.19+ or 22.12+ (Node.js 18 support dropped)
- **Release Cycle**: Regular minor releases with performance improvements

### Official Documentation
- **Main Documentation**: https://vite.dev/
- **Getting Started**: https://vite.dev/guide/
- **Configuration Reference**: https://vite.dev/config/
- **Plugin Development**: https://vite.dev/guide/api-plugin
- **Migration Guide**: https://vite.dev/guide/migration
- **GitHub Repository**: https://github.com/vitejs/vite

### Key Features for RealWorld
- **Lightning Fast HMR**: Sub-second hot module replacement
- **Native ES Modules**: No bundling during development
- **TypeScript Support**: Built-in TypeScript compilation
- **Tree Shaking**: Efficient dead code elimination
- **Code Splitting**: Automatic chunking and lazy loading
- **Asset Optimization**: Image compression, WebP support
- **CSS Features**: PostCSS, CSS modules, preprocessors

### Installation & Setup
```bash
# Create new Vite React TypeScript project
npm create vite@latest realworld-frontend -- --template react-ts

# Or install in existing project
npm install -D vite@^7.1.4
npm install -D @vitejs/plugin-react@^4.0.0
```

### Vite Configuration for RealWorld
```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [
    react({
      // Enable React Fast Refresh
      fastRefresh: true,
      // JSX automatic runtime
      jsxRuntime: 'automatic'
    })
  ],
  
  // Path resolution for clean imports
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/components': path.resolve(__dirname, './src/components'),
      '@/pages': path.resolve(__dirname, './src/pages'),
      '@/lib': path.resolve(__dirname, './src/lib'),
      '@/hooks': path.resolve(__dirname, './src/hooks'),
      '@/types': path.resolve(__dirname, './src/types')
    }
  },

  // Development server configuration
  server: {
    port: 3000,
    host: true, // Listen on all addresses
    proxy: {
      // Proxy API calls to backend during development
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  },

  // Build optimization
  build: {
    target: 'es2022',
    outDir: 'dist',
    sourcemap: true,
    minify: 'esbuild',
    
    // Chunk splitting strategy
    rollupOptions: {
      output: {
        manualChunks: {
          // Separate vendor chunks
          'vendor': ['react', 'react-dom'],
          'router': ['@tanstack/react-router'],
          'query': ['@tanstack/react-query'],
          'ui': ['lucide-react', '@radix-ui/react-slot']
        }
      }
    },
    
    // Performance budgets
    chunkSizeWarningLimit: 1000
  },

  // Environment variable prefix
  envPrefix: 'VITE_',

  // CSS configuration
  css: {
    postcss: './postcss.config.js',
    devSourcemap: true
  },

  // Optimization
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      '@tanstack/react-query',
      '@tanstack/react-router'
    ]
  }
})
```

### Environment Variables
```bash
# .env.development
VITE_API_URL=http://localhost:8080
VITE_APP_TITLE=RealWorld
VITE_DEBUG=true

# .env.production
VITE_API_URL=https://api.realworld-app.com
VITE_APP_TITLE=RealWorld
VITE_DEBUG=false
```

### Build Scripts
```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "build:analyze": "vite build --mode analyze",
    "build:production": "NODE_ENV=production vite build"
  }
}
```

## Performance Features

### Development Performance
- **Cold Start**: ~200ms initial startup time
- **HMR**: Sub-50ms hot module replacement
- **File Watching**: Efficient file system monitoring
- **Memory Usage**: Optimized for large codebases

### Production Optimizations
- **Bundle Splitting**: Automatic vendor and route-based chunks
- **Tree Shaking**: Dead code elimination via Rollup
- **Minification**: ESBuild-powered JavaScript/CSS minification
- **Asset Optimization**: Image compression, format conversion
- **Preload Directives**: Automatic resource hints generation

### Build Output Structure
```
dist/
├── assets/
│   ├── index-[hash].js          # Main application bundle
│   ├── vendor-[hash].js         # Third-party libraries
│   ├── router-[hash].js         # Routing logic
│   ├── ui-[hash].js            # UI component library
│   └── index-[hash].css        # Compiled styles
├── index.html                   # Entry point with asset links
└── vite.svg                     # Static assets
```

## Plugin Ecosystem

### Core Plugins for RealWorld
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [
    // React support with Fast Refresh
    react(),
    
    // Bundle analyzer (development)
    process.env.ANALYZE && bundleAnalyzer(),
    
    // PWA support (future enhancement)
    // VitePWA({ registerType: 'autoUpdate' })
  ]
})
```

### Available Plugins
- **@vitejs/plugin-react**: React Fast Refresh and JSX support
- **vite-bundle-analyzer**: Bundle size analysis
- **vite-plugin-pwa**: Progressive Web App features
- **vite-plugin-eslint**: ESLint integration
- **vite-plugin-windicss**: Alternative CSS framework support

## Integration with RealWorld Stack

### TypeScript Integration
- **Zero Configuration**: Built-in TypeScript support
- **Type Checking**: Parallel type checking with tsc
- **Path Mapping**: Automatic alias resolution
- **Declaration Files**: .d.ts generation for libraries

### Tailwind CSS Integration
```typescript
// vite.config.ts - CSS configuration
export default defineConfig({
  css: {
    postcss: {
      plugins: [
        require('tailwindcss'),
        require('autoprefixer'),
      ]
    }
  }
})
```

### Testing Integration
- **Vitest**: Vite-powered testing (compatible with Vite config)
- **Coverage**: Built-in c8 coverage reports
- **Browser Mode**: In-browser testing support

## Migration from Webpack/CRA

### Key Differences
- **No Webpack Config**: Minimal configuration required
- **ES Modules**: Native ESM during development
- **Faster Builds**: 10-100x faster than traditional bundlers
- **Better DX**: Superior developer experience

### Migration Steps
1. Install Vite and React plugin
2. Update build scripts in package.json
3. Convert webpack.config.js to vite.config.ts
4. Update environment variable prefixes (REACT_APP_ → VITE_)
5. Test build output and deployment process

This Vite configuration provides the fastest possible development experience with production-ready build optimization for the RealWorld application.