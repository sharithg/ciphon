import { createFileRoute } from '@tanstack/react-router'
import GithubLogin from '../components/github-login'

export const Route = createFileRoute('/login')({
  component: GithubLogin,
})
