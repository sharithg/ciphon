import { createLazyFileRoute } from "@tanstack/react-router";
import { useEffect, useMemo, useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { useGetSteps } from "@/hooks/react-query/use-workflows";
import CommandOutput from "@/components/command-output";

export const Route = createLazyFileRoute(
  "/pipelines/workflows/$workflowId/jobs/$jobId"
)({
  component: () => (
    <>
      <JobSteps />
    </>
  ),
});

function JobSteps() {
  const [expandedCommands, setExpandedCommands] = useState<string[]>([]);

  const [userExpandedCommands, setUserExpandedCommands] = useState<string[]>(
    []
  );

  const { workflowId, jobId } = Route.useParams();

  const { data } = useGetSteps(workflowId, jobId);

  const steps = useMemo(() => data?.data ?? [], [data?.data]);

  const toggleUserExpand = (id: string) => {
    setUserExpandedCommands((prev) =>
      prev.includes(id)
        ? prev.filter((commandId) => commandId !== id)
        : [...prev, id]
    );
  };

  useEffect(() => {
    if (steps.length) {
      const runningCommands = steps
        .filter(
          (s) => s.status === "running" && !userExpandedCommands.includes(s.id)
        )
        .map((s) => s.id);

      setExpandedCommands(() => [...userExpandedCommands, ...runningCommands]);
    }
  }, [steps, userExpandedCommands]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case "success":
        return "bg-green-500";
      case "failed":
        return "bg-red-500";
      case "running":
        return "bg-yellow-500";
      default:
        return "bg-gray-500";
    }
  };

  return (
    <div className="flex-1 overflow-hidden">
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-4">Command Runs</h1>
        {steps.map((command) => (
          <div key={command.id} className="mb-4">
            <div
              className="flex items-center justify-between bg-muted/60 p-3 rounded-t-lg cursor-pointer"
              onClick={() => toggleUserExpand(command.id)}
            >
              <div className="flex items-center space-x-3">
                {expandedCommands.includes(command.id) ? (
                  <ChevronDown className="h-5 w-5" />
                ) : (
                  <ChevronRight className="h-5 w-5" />
                )}
                <span>{command.name || command.type}</span>
              </div>
              <Badge
                className={`${command.status ? getStatusColor(command.status) : "bg-black-500"} text-white`}
              >
                {command.status}
              </Badge>
            </div>
            {expandedCommands.includes(command.id) && (
              <CommandOutput
                workflowId={workflowId}
                jobId={jobId}
                stepId={command.id}
                status={command.status}
              />
            )}
            {!expandedCommands.includes(command.id) && <Separator />}
          </div>
        ))}
      </div>
    </div>
  );
}
