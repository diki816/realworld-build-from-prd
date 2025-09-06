import { createFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/register')({
  component: Register,
})

function Register() {
  return (
    <div className="container mx-auto py-12 px-4 max-w-md">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl text-gray-800">Sign up</CardTitle>
          <Link to="/login" className="text-conduit-green hover:underline">
            Have an account?
          </Link>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <form className="space-y-4">
            <div>
              <Input
                type="text"
                placeholder="Username"
                className="form-control form-control-lg"
              />
            </div>
            
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
              Sign up
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}