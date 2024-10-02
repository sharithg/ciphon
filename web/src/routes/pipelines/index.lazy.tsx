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

const StatusBadge = ({ status }: { status: PipelineRun["status"] }) => {
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
  };

  const { label, icon: Icon, className } = statusConfig[status];

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
  const [selectedProject, setSelectedProject] = useState<string | undefined>();
  const [selectedBranch, setSelectedBranch] = useState<string | undefined>();
  const [selectedDate, setSelectedDate] = useState<Date | undefined>();

  const projects = Array.from(new Set(pipelineRuns.map((run) => run.project)));
  const branches = Array.from(new Set(pipelineRuns.map((run) => run.branch)));

  const filteredRuns = pipelineRuns.filter(
    (run) =>
      (!selectedProject || run.project === selectedProject) &&
      (!selectedBranch || run.branch === selectedBranch) &&
      (!selectedDate ||
        run.timestamp.startsWith(format(selectedDate, "yyyy-MM-dd")))
  );

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
          <TableRow>
            <TableHead>Project</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Workflow</TableHead>
            <TableHead>Trigger</TableHead>
            <TableHead>Timestamp</TableHead>
            <TableHead>Duration</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredRuns.map((run) => (
            <TableRow key={run.id}>
              <TableCell className="font-medium">{run.project}</TableCell>
              <TableCell>
                <StatusBadge status={run.status} />
              </TableCell>
              <TableCell>{run.workflow}</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <GitBranchIcon className="w-4 h-4" />
                  <span>{run.branch}</span>
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-500">
                  <GitCommitIcon className="w-4 h-4" />
                  <span>{run.commit}</span>
                </div>
              </TableCell>
              <TableCell>{run.timestamp}</TableCell>
              <TableCell>
                <div className="flex items-center gap-1">
                  <ClockIcon className="w-4 h-4" />
                  <span>{run.duration}</span>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
