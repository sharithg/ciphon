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
import { formatDistance } from "date-fns";

export const Route = createLazyFileRoute("/projects")({
  component: Projects,
});

function Projects() {
  const [searchTerm, setSearchTerm] = useState("");

  const repos = useGetRepos();

  console.log(repos.data?.data);

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
        {repos.data?.data.map((project) => (
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
                Last updated:{" "}
                {formatDistance(new Date(project.lastUpdated), new Date(), {
                  addSuffix: true,
                })}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
      {repos.data?.data.length === 0 && (
        <p className="text-center text-muted-foreground">No projects found.</p>
      )}
    </div>
  );
}
