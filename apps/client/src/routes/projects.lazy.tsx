import React from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  GitHubLogoIcon,
  MagnifyingGlassIcon,
  PlusIcon,
} from "@radix-ui/react-icons";
import { useGetRepos } from "../hooks/use-github";

type Project = {
  id: string;
  name: string;
  description: string;
  lastUpdated: string;
};

const mockProjects: Project[] = [
  {
    id: "1",
    name: "awesome-project",
    description: "An awesome project",
    lastUpdated: "2 days ago",
  },
  {
    id: "2",
    name: "cool-app",
    description: "A cool application",
    lastUpdated: "5 hours ago",
  },
  {
    id: "3",
    name: "my-website",
    description: "Personal website",
    lastUpdated: "1 week ago",
  },
];
export const Route = createLazyFileRoute("/projects")({
  component: Projects,
});

function Projects() {
  const [projects, setProjects] = useState<Project[]>(mockProjects);
  const [searchTerm, setSearchTerm] = useState("");

  const repos = useGetRepos();

  console.log(repos.data?.data);

  const filteredProjects = projects.filter((project) =>
    project.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleConnectNewProject = () => {
    // Implement the logic to connect a new project
    console.log("Connecting new project...");
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Projects</h2>
        <Button onClick={handleConnectNewProject}>
          <PlusIcon className="mr-2 h-4 w-4" /> Connect New Project
        </Button>
      </div>
      <div className="relative">
        <MagnifyingGlassIcon className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search projects..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="pl-8"
        />
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {filteredProjects.map((project) => (
          <Card key={project.id}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {project.name}
              </CardTitle>
              <GitHubLogoIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <CardDescription>{project.description}</CardDescription>
              <p className="text-xs text-muted-foreground mt-2">
                Last updated: {project.lastUpdated}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
      {filteredProjects.length === 0 && (
        <p className="text-center text-muted-foreground">No projects found.</p>
      )}
    </div>
  );
}
