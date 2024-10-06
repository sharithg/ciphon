import { createLazyFileRoute } from "@tanstack/react-router";
import Pipelines from "../../components/pipelines";
export const Route = createLazyFileRoute("/pipelines/")({
  component: PipelinesProvider,
});

function PipelinesProvider() {
  return (
    <>
      <Pipelines />
    </>
  );
}
