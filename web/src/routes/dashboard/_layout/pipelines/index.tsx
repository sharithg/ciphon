import { createFileRoute } from "@tanstack/react-router";
import Pipelines from "@/components/pipelines";

export const Route = createFileRoute("/dashboard/_layout/pipelines/")({
  component: PipelinesProvider,
});

function PipelinesProvider() {
  return (
    <>
      <Pipelines />
    </>
  );
}
