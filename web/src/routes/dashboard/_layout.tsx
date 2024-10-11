import { createFileRoute, redirect } from "@tanstack/react-router";
import { Outlet } from "@tanstack/react-router";
import { ContentLayout } from "@/components/admin-panel/content-layout";
import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import { Toaster } from "@/components/ui/toaster";
import WorkflowBreadcrumbs from "@/components/workflow-breadcrumbs";
import { isAuthenticated } from "@/hooks/user-auth";

export const Route = createFileRoute("/dashboard/_layout")({
  component: Root,
  beforeLoad: async ({ location }) => {
    const isAuthed = await isAuthenticated();
    if (!isAuthed) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
});

function Root() {
  console.log("root!");
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
