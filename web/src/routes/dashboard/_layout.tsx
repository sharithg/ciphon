import { createFileRoute } from "@tanstack/react-router";
import { Outlet } from "@tanstack/react-router";
import { ContentLayout } from "@/components/admin-panel/content-layout";
import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import { Toaster } from "@/components/ui/toaster";
import WorkflowBreadcrumbs from "@/components/workflow-breadcrumbs";

export const Route = createFileRoute("/dashboard/_layout")({
  component: Root,
});

function Root() {
  console.log("rooot");
  return (
    <>
      <AdminPanelLayout>
        <ContentLayout title="Account">
          <WorkflowBreadcrumbs>
            <Outlet />
          </WorkflowBreadcrumbs>
          <Toaster />
        </ContentLayout>
      </AdminPanelLayout>
    </>
  );
}
