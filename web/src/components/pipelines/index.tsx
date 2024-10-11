import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  GitBranchIcon,
  GitCommitIcon,
  CheckCircleIcon,
  XCircleIcon,
  AlertCircleIcon,
  ClockIcon,
  Play,
  Loader2,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { format, formatDistance } from "date-fns";
import {
  useGetWorkflows,
  useRunWorkflow,
} from "@/hooks/react-query/use-workflows";
import { Link } from "@tanstack/react-router";
import { cn } from "../../lib/utils";

const StatusBadge = ({ status }: { status: string }) => {
  const statusConfig = {
    success: {
      label: "Success",
      icon: CheckCircleIcon,
      className: "bg-green-500",
    },
    failed: { label: "Failed", icon: XCircleIcon, className: "bg-red-500" },
    running: {
      label: "Running",
      icon: () => <Loader2 className="animate-spin" />,
      className: "bg-blue-500",
    },
    canceled: {
      label: "Canceled",
      icon: AlertCircleIcon,
      className: "bg-yellow-500",
    },
    not_started: {
      label: "Not Started",
      icon: AlertCircleIcon,
      className: "bg-gray-500",
    },
  };

  let statusKey = status as keyof typeof statusConfig;
  if (!status) {
    statusKey = "not_started";
  }

  const { label, icon: Icon, className } = statusConfig[statusKey];

  return (
    <Badge
      variant="secondary"
      className={`${className} text-white flex items-center gap-1`}
    >
      <Icon className="w-3 h-3" />
      {label}
    </Badge>
  );
};

const Pipelines = () => {
  const [, setSelectedProject] = useState<string | undefined>();
  const [, setSelectedBranch] = useState<string | undefined>();
  const [selectedDate, setSelectedDate] = useState<Date | undefined>();
  const { data, refetch } = useGetWorkflows();
  const mutation = useRunWorkflow();

  // useWorkflowEvents(`${API_URL}/sse/workflows/run-events`, (event) => {
  //   if (event.type === "workflow") {
  //     refetch();
  //   }
  // });

  return (
    <div className="container mx-auto py-10">
      <h1 className="text-2xl font-bold mb-4">Pipeline Runs</h1>
      <div className="flex flex-wrap gap-4 mb-6">
        <Select onValueChange={setSelectedProject}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select Project" />
          </SelectTrigger>
          <SelectContent>
            {/* {projects.map((project) => (
              <SelectItem key={project} value={project}>
                {project}
              </SelectItem>
            ))} */}
          </SelectContent>
        </Select>
        <Select onValueChange={setSelectedBranch}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select Branch" />
          </SelectTrigger>
          <SelectContent>
            {/* {branches.map((branch) => (
              <SelectItem key={branch} value={branch}>
                {branch}
              </SelectItem>
            ))} */}
          </SelectContent>
        </Select>
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline">
              {selectedDate ? format(selectedDate, "PPP") : "Pick a date"}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0">
            <Calendar
              mode="single"
              selected={selectedDate}
              onSelect={setSelectedDate}
              initialFocus
            />
          </PopoverContent>
        </Popover>
        <Button
          onClick={() => {
            setSelectedProject(undefined);
            setSelectedBranch(undefined);
            setSelectedDate(undefined);
          }}
        >
          Clear Filters
        </Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow className="h-8">
            <TableHead className="p-2 text-xs">Project</TableHead>
            <TableHead className="p-2 text-xs">Status</TableHead>
            <TableHead className="p-2 text-xs">Workflow</TableHead>
            <TableHead className="p-2 text-xs">Trigger</TableHead>
            <TableHead className="p-2 text-xs">Timestamp</TableHead>
            <TableHead className="p-2 text-xs">Duration</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {(data?.data ?? []).map((run) => (
            <TableRow key={run.workflowId} className="h-8">
              <TableCell className="p-2 text-sm font-medium">
                {run.repoName}
              </TableCell>
              <TableCell className="p-2 text-sm">
                <StatusBadge status={run.status} />
              </TableCell>
              <TableCell className="p-2 text-md">
                <Link
                  disabled={!run.status}
                  to={`/dashboard/pipelines/workflows/${run.workflowId}`}
                  className={cn(
                    run.status ? "text-blue-500 hover:underline" : ""
                  )}
                >
                  {run.workflowName}
                </Link>
              </TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <GitBranchIcon className="w-3 h-3" />
                  <span>{run.branch}</span>
                </div>
                <div className="flex items-center gap-1 text-xs text-gray-500">
                  <GitCommitIcon className="w-3 h-3" />
                  <span>{(run.commitSha ?? "").slice(0, 7)}</span>
                </div>
              </TableCell>
              <TableCell className="p-2 text-sm">
                {format(new Date(run.createdAt), "MMMM do, yyyy 'at' hh:mm a")}
              </TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <ClockIcon className="w-3 h-3" />
                  <span>
                    {run.duration
                      ? `${formatDistance(0, run.duration * 1000, { includeSeconds: true })}`
                      : ""}
                  </span>
                </div>
              </TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <Button
                    variant="outline"
                    size="icon"
                    disabled={run.status === "running" || mutation.isLoading}
                    onClick={async () => {
                      await mutation.mutateAsync(run.workflowId);
                      refetch();
                    }}
                  >
                    <Play className="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default Pipelines;
