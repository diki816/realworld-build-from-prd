import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

export const Route = createRootRoute({
  component: () => (
    <>
      <div className="min-h-screen bg-white">
        {/* Navigation Header */}
        <nav className="bg-white border-b border-gray-200 px-4 py-2">
          <div className="container mx-auto flex justify-between items-center">
            <Link to="/" className="text-2xl font-bold text-conduit-green">
              conduit
            </Link>
            <div className="flex space-x-4">
              <Link to="/" className="text-gray-700 hover:text-conduit-green">
                Home
              </Link>
              <Link to="/login" className="text-gray-700 hover:text-conduit-green">
                Sign in
              </Link>
              <Link to="/register" className="text-gray-700 hover:text-conduit-green">
                Sign up
              </Link>
            </div>
          </div>
        </nav>
        
        {/* Main Content */}
        <main>
          <Outlet />
        </main>
        
        {/* Footer */}
        <footer className="bg-gray-100 py-4 mt-8">
          <div className="container mx-auto text-center text-gray-600">
            <p>© 2025 Conduit. An interactive learning project from Thinkster.</p>
          </div>
        </footer>
      </div>
      <TanStackRouterDevtools />
    </>
  ),
})