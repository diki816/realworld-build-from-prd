import { createFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
  return (
    <div className="container mx-auto py-12 px-4 max-w-md">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl text-gray-800">Sign in</CardTitle>
          <Link to="/register" className="text-conduit-green hover:underline">
            Need an account?
          </Link>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <form className="space-y-4">
            <div>
              <Input
                type="email"
                placeholder="Email"
                className="form-control form-control-lg"
              />
            </div>
            
            <div>
              <Input
                type="password"
                placeholder="Password"
                className="form-control form-control-lg"
              />
            </div>
            
            <Button
              type="submit"
              className="btn btn-primary btn-lg w-full"
              size="lg"
            >
              Sign in
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}