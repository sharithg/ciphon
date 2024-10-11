import { createFileRoute, FileRoutesByPath } from "@tanstack/react-router";
import { useEffect } from "react";
import { useAuth } from "@/hooks/user-auth";

export const Route = createFileRoute("/login/github/callback")({
  component: Callback,
});

function Callback() {
  const { handleCallback } = useAuth();

  useEffect(() => {
    handleCallback()
      .then(() => {
        (window.location
          .href as FileRoutesByPath[keyof FileRoutesByPath]["fullPath"]) =
          "/dashboard/pipelines";
      })
      .catch(() => {
        (window.location
          .href as FileRoutesByPath[keyof FileRoutesByPath]["fullPath"]) =
          "/login";
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return <></>;
}
