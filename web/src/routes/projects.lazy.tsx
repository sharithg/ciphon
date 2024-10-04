import React from "react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { GitHubLogoIcon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { useGetRepos } from "../hooks/react-query/use-github";
import { formatDistance } from "date-fns";
import ConnectRepo from "../components/connect-repo";

export const Route = createLazyFileRoute("/projects")({
  component: Projects,
});

function Projects() {
  const [searchTerm, setSearchTerm] = useState("");

  const repos = useGetRepos();

  const repoData = repos.data?.data ?? [];

  return (
    <div className="space-y-4 pt-5">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Projects</h2>
        <ConnectRepo />
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
        {repoData.map((project) => (
          <Card key={project.repoId}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {project.name}
              </CardTitle>
              <GitHubLogoIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <CardDescription>{project.description}</CardDescription>
              <p className="text-xs text-muted-foreground mt-2">
                Created:{" "}
                {formatDistance(new Date(project.repoCreatedAt), new Date(), {
                  addSuffix: true,
                })}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
      {repoData.length === 0 && (
        <p className="text-center text-muted-foreground">No projects found.</p>
      )}
    </div>
  );
}
