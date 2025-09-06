import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div>
      {/* Hero Section */}
      <div className="bg-conduit-green text-white text-center py-20">
        <div className="container mx-auto">
          <h1 className="text-4xl font-bold mb-2">conduit</h1>
          <p className="text-xl">A place to share your knowledge.</p>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto py-8 px-4">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Articles Feed */}
          <div className="lg:col-span-2">
            <div className="border-b border-gray-200 mb-4">
              <div className="flex space-x-8">
                <button className="pb-2 border-b-2 border-conduit-green text-conduit-green font-medium">
                  Global Feed
                </button>
              </div>
            </div>
            
            {/* Article List Placeholder */}
            <div className="space-y-4">
              <div className="border-t border-gray-200 pt-4">
                <div className="flex items-center space-x-2 text-sm text-gray-600 mb-2">
                  <img
                    src="https://via.placeholder.com/32"
                    alt="Author"
                    className="w-8 h-8 rounded-full"
                  />
                  <span>Loading articles...</span>
                </div>
                <h2 className="text-xl font-semibold text-gray-800 mb-2">
                  Welcome to Conduit
                </h2>
                <p className="text-gray-600 mb-2">
                  This is your global feed. Articles will appear here once the backend is connected.
                </p>
                <div className="flex items-center space-x-2 text-sm text-gray-500">
                  <span>January 1, 2025</span>
                </div>
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <div className="bg-conduit-gray-light p-4 rounded">
              <h3 className="font-semibold text-gray-800 mb-3">Popular Tags</h3>
              <div className="flex flex-wrap gap-1">
                <span className="tag-default">welcome</span>
                <span className="tag-default">introduction</span>
                <span className="tag-default">getting-started</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}