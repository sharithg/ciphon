import { createLazyFileRoute } from "@tanstack/react-router";
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
  PlayCircleIcon,
  CheckCircleIcon,
  XCircleIcon,
  AlertCircleIcon,
  ClockIcon,
  Play,
} from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
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
import { format } from "date-fns";
import { useGetWorkflows } from "@/hooks/react-query/use-workflows";

export const Route = createLazyFileRoute("/pipelines/")({
  component: Pipelines,
});

type PipelineRun = {
  id: string;
  project: string;
  status: "success" | "failed" | "running" | "canceled";
  workflow: string;
  branch: string;
  commit: string;
  timestamp: string;
  duration: string;
};

const pipelineRuns: PipelineRun[] = [
  {
    id: "1",
    project: "Frontend App",
    status: "success",
    workflow: "Build and Test",
    branch: "main",
    commit: "a1b2c3d",
    timestamp: "2023-09-30 14:30:00",
    duration: "3m 45s",
  },
  {
    id: "2",
    project: "Backend API",
    status: "failed",
    workflow: "Integration Tests",
    branch: "feature/auth",
    commit: "e4f5g6h",
    timestamp: "2023-09-30 13:15:00",
    duration: "5m 20s",
  },
  {
    id: "3",
    project: "Mobile App",
    status: "running",
    workflow: "Build iOS",
    branch: "release/v2.0",
    commit: "i7j8k9l",
    timestamp: "2023-09-30 15:00:00",
    duration: "2m 10s",
  },
  {
    id: "4",
    project: "Data Pipeline",
    status: "canceled",
    workflow: "ETL Process",
    branch: "fix/data-import",
    commit: "m1n2o3p",
    timestamp: "2023-09-30 12:45:00",
    duration: "1m 55s",
  },
];

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
      icon: PlayCircleIcon,
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

function Pipelines() {
  const [, setSelectedProject] = useState<string | undefined>();
  const [, setSelectedBranch] = useState<string | undefined>();
  const [selectedDate, setSelectedDate] = useState<Date | undefined>();
  const { data } = useGetWorkflows();

  const projects = Array.from(new Set(pipelineRuns.map((run) => run.project)));
  const branches = Array.from(new Set(pipelineRuns.map((run) => run.branch)));

  // const filteredRuns = pipelineRuns.filter(
  //   (run) =>
  //     (!selectedProject || run.project === selectedProject) &&
  //     (!selectedBranch || run.branch === selectedBranch) &&
  //     (!selectedDate ||
  //       run.timestamp.startsWith(format(selectedDate, "yyyy-MM-dd")))
  // );

  return (
    <div className="container mx-auto py-10">
      <h1 className="text-2xl font-bold mb-4">Pipeline Runs</h1>
      <div className="flex flex-wrap gap-4 mb-6">
        <Select onValueChange={setSelectedProject}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select Project" />
          </SelectTrigger>
          <SelectContent>
            {projects.map((project) => (
              <SelectItem key={project} value={project}>
                {project}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select onValueChange={setSelectedBranch}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Select Branch" />
          </SelectTrigger>
          <SelectContent>
            {branches.map((branch) => (
              <SelectItem key={branch} value={branch}>
                {branch}
              </SelectItem>
            ))}
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
            {" "}
            {/* Reduce height */}
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
              {" "}
              {/* Reduce height */}
              <TableCell className="p-2 text-sm font-medium">
                {run.repoName}
              </TableCell>
              <TableCell className="p-2 text-sm">
                <StatusBadge status={run.status} />
              </TableCell>
              <TableCell className="p-2 text-sm">{run.workflowName}</TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <GitBranchIcon className="w-3 h-3" />
                  <span>{run.branch}</span>
                </div>
                <div className="flex items-center gap-1 text-xs text-gray-500">
                  <GitCommitIcon className="w-3 h-3" />
                  <span>{run.commitSha.slice(0, 7)}</span>
                </div>
              </TableCell>
              <TableCell className="p-2 text-sm">
                {format(new Date(run.createdAt), "MMMM do, yyyy 'at' hh:mm a")}
              </TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <ClockIcon className="w-3 h-3" />
                  <span>{run.duration}</span>
                </div>
              </TableCell>
              <TableCell className="p-2 text-sm">
                <div className="flex items-center gap-1">
                  <Button variant="outline" size="icon">
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
}
