"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { GitBranch } from "lucide-react";
import {
  TConnectRepo,
  useConnectRepo,
  useGetNewRepos,
} from "../hooks/react-query/use-github";

export default function ConnectRepo() {
  const [isOpen, setIsOpen] = useState(false);
  const repos = useGetNewRepos();

  const mutation = useConnectRepo();
  const repoData = repos.data?.data ?? [];

  const handleConnect = async (repoData: TConnectRepo) => {
    mutation.mutate(repoData);
  };

  return (
    <>
      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogTrigger asChild>
          <Button>
            <GitBranch className="mr-2 h-4 w-4" />
            Connect Repo
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Connect to a Repository</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            {repoData.length ? (
              repoData.map((repo) => (
                <Button
                  key={repo.repoId}
                  onClick={() =>
                    handleConnect({
                      name: repo.name,
                      owner: repo.owner,
                    })
                  }
                  disabled={mutation.isLoading}
                  className="justify-start"
                >
                  <GitBranch className="mr-2 h-4 w-4" />
                  {repo.name}
                </Button>
              ))
            ) : (
              <h1>No new repos to connect.</h1>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
