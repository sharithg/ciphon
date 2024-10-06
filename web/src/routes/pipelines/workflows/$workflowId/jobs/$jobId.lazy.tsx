import { createLazyFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { useGetSteps } from "../../../../../hooks/react-query/use-workflows";
import useWorkflowEvents from "../../../../../hooks/react-query/use-sse";

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

  const { workflowId, jobId } = Route.useParams();

  const { data, refetch } = useGetSteps(workflowId, jobId);
  const toggleExpand = (id: string) => {
    setExpandedCommands((prev) =>
      prev.includes(id)
        ? prev.filter((commandId) => commandId !== id)
        : [...prev, id]
    );
  };

  useWorkflowEvents((event) => {
    if (event.type === "step") {
      refetch();
    }
  });

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
        <ScrollArea className="h-[calc(100vh-8rem)] pr-4">
          {(data?.data ?? []).map((command) => (
            <div key={command.id} className="mb-4">
              <div
                className="flex items-center justify-between bg-muted/60 p-3 rounded-t-lg cursor-pointer"
                onClick={() => toggleExpand(command.id)}
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
                <div className="bg-muted/30 p-3 rounded-b-lg">
                  <pre className="whitespace-pre-wrap font-mono text-sm">
                    echo test
                  </pre>
                </div>
              )}
              {!expandedCommands.includes(command.id) && <Separator />}
            </div>
          ))}
        </ScrollArea>
      </div>
    </div>
  );
}
