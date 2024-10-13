import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Link, useRouterState } from "@tanstack/react-router";
import { useAtom } from "jotai";
import { jobs, workflows } from "./atoms/workflows";
import { useMemo } from "react";

function getPathType(path: string) {
  const workflowRegex = /^\/dashboard\/pipelines\/workflows\/([a-f0-9-]+)$/;
  const jobsRegex =
    /^\/dashboard\/pipelines\/workflows\/([a-f0-9-]+)\/jobs\/([a-f0-9-]+)$/;

  const jobsMatch = jobsRegex.exec(path);
  if (jobsMatch) {
    return {
      type: "jobs",
      workflowId: jobsMatch[1],
      jobId: jobsMatch[2],
    };
  }

  const workflowMatch = workflowRegex.exec(path);
  if (workflowMatch) {
    return {
      type: "workflow",
      workflowId: workflowMatch[1],
    };
  }

  return null;
}

const WorkflowBreadcrumbs = ({ children }: { children: React.ReactNode }) => {
  const [allJobs] = useAtom(jobs);
  const [allWorkflows] = useAtom(workflows);
  const router = useRouterState();
  const currentPath = getPathType(router.location.pathname);

  const selectedWorkflow = useMemo(
    () =>
      (allWorkflows ?? []).find(
        (j) => j.workflowId === currentPath?.workflowId
      ),
    [allWorkflows, currentPath]
  );
  const selectedJob = useMemo(
    () => (allJobs ?? []).find((j) => j.id === currentPath?.jobId),
    [allJobs, currentPath]
  );

  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          {(currentPath?.type === "workflow" ||
            currentPath?.type === "jobs") && (
            <>
              <BreadcrumbItem>
                <BreadcrumbLink>
                  <Link to="/dashboard/pipelines">Workflows</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>
                  {selectedWorkflow?.workflowName}
                </BreadcrumbLink>
              </BreadcrumbItem>
            </>
          )}
          {selectedJob && currentPath?.type === "jobs" && (
            <>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>
                  <Link
                    to={`/dashboard/pipelines/workflows/${currentPath.workflowId}`}
                  >
                    Jobs
                  </Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>{selectedJob.name}</BreadcrumbLink>
              </BreadcrumbItem>
            </>
          )}
        </BreadcrumbList>
      </Breadcrumb>
      {children}
    </>
  );
};

export default WorkflowBreadcrumbs;
