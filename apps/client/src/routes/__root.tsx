import { createRootRoute, Outlet } from "@tanstack/react-router";
import { ThemeProvider } from "@/components/theme-provider";
import { ContentLayout } from "../components/admin-panel/content-layout";
import AdminPanelLayout from "../components/admin-panel/admin-panel-layout";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "../components/ui/toaster";

export const Route = createRootRoute({
  component: Root,
});

const queryClient = new QueryClient();

function Root() {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
          <AdminPanelLayout>
            <ContentLayout title="Account">
              <Outlet />
              <Toaster />
            </ContentLayout>
          </AdminPanelLayout>
        </ThemeProvider>
      </QueryClientProvider>
    </>
  );
}
