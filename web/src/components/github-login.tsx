import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Github } from "lucide-react";
import { API_URL } from "../hooks/react-query/constants";

export default function GithubLogin() {
  const handleGitHubLogin = () => {
    window.location.href = `${API_URL}/auth/login/github`;
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-zinc-800">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            Login
          </CardTitle>
          <CardDescription className="text-center">
            Sign in to your account using GitHub
          </CardDescription>
        </CardHeader>
        <CardContent className="flex justify-center">
          <Button className="w-full max-w-sm" onClick={handleGitHubLogin}>
            <Github className="mr-2 h-4 w-4" />
            Sign in with GitHub
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
