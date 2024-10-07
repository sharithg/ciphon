import React, { useEffect, useMemo, useRef } from "react";
import { useGetCommandOutput } from "@/hooks/react-query/use-workflows";
import { API_URL } from "../hooks/react-query/constants";
import useWorkflowEvents from "../hooks/react-query/use-sse";

const CommandOutput = (input: {
  workflowId: string;
  jobId: string;
  stepId: string;
}) => {
  const { data, refetch } = useGetCommandOutput(
    input.workflowId,
    input.jobId,
    input.stepId
  );
  const containerRef = useRef<HTMLDivElement | null>(null);
  const outputs = useMemo(() => data?.data ?? [], [data?.data]);

  useEffect(() => {
    // Scroll to the bottom of the container when outputs change
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [outputs]);

  useWorkflowEvents(
    `${API_URL}/sse/steps/run-events/${input.stepId}`,
    (event) => {
      if (event.type === "step") {
        refetch();
      }
    }
  );

  return (
    <div
      ref={containerRef}
      className="bg-muted/30 p-3 rounded-b-lg overflow-y-auto max-h-96"
    >
      <pre className="whitespace-pre-wrap font-mono text-sm">
        {outputs.map((v, index) => (
          <span key={index}>{v.stdout}</span>
        ))}
      </pre>
    </div>
  );
};

export default CommandOutput;
