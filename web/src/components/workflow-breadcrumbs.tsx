import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Link, useRouterState } from "@tanstack/react-router";
import { useAtom } from "jotai";
import { selectedJobAtom, selectedWorkflowAtom } from "./atoms/workflows";

function getPathType(path: string) {
  const workflowRegex = /^\/pipelines\/workflows\/[a-f0-9-]+$/;
  const jobsRegex = /^\/pipelines\/workflows\/[a-f0-9-]+\/jobs\/[a-f0-9-]+$/;

  if (jobsRegex.test(path)) {
    return "jobs";
  } else if (workflowRegex.test(path)) {
    return "workflow";
  }
  return null;
}

const WorkflowBreadcrumbs = ({ children }: { children: React.ReactNode }) => {
  const [selectedJob] = useAtom(selectedJobAtom);
  const [selectedWorkflow] = useAtom(selectedWorkflowAtom);
  const router = useRouterState();
  const currentPath = getPathType(router.location.pathname);

  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          {selectedWorkflow &&
            (currentPath === "workflow" || currentPath === "jobs") && (
              <>
                <BreadcrumbItem>
                  <BreadcrumbLink>
                    <Link to="/pipelines">Workflows</Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  <BreadcrumbLink>{selectedWorkflow.name}</BreadcrumbLink>
                </BreadcrumbItem>
              </>
            )}
          {selectedJob && selectedWorkflow && currentPath === "jobs" && (
            <>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>
                  <Link to={`/pipelines/workflows/${selectedWorkflow.id}`}>
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
