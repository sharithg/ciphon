import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  ReactFlow,
  Controls,
  Background,
  Node,
  Position,
  Edge,
} from "@xyflow/react";

import "@xyflow/react/dist/base.css";
import CustomNode from "@/components/custom-node";
import { useGetJobs } from "@/hooks/react-query/use-workflows";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/pipelines/workflows/$workflowId/")({
  component: () => (
    <>
      <Pipeline />
    </>
  ),
});

const nodeTypes = {
  customNode: CustomNode,
};

const edges: Edge[] = [];

function Pipeline() {
  const params = Route.useParams();
  const { data } = useGetJobs(params.workflowId);

  const [nodes, setNodes] = useState<Node[]>([]);

  const nav = useNavigate();

  // useWorkflowEvents(`${API_URL}/sse/workflows/run-events`, (event) => {
  //   if (event.type === "job") {
  //     refetch();
  //   }
  // });

  useEffect(() => {
    if (data?.data) {
      setNodes(
        data.data.map((v) => ({
          id: v.id,
          sourcePosition: Position.Right,
          type: "customNode",
          data: { label: v.name, status: v.status },
          position: { x: 0, y: 80 },
        }))
      );
    }
  }, [data]);

  return (
    <>
      <div style={{ height: 500, width: "100%" }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          fitView
          nodesConnectable={false}
          elementsSelectable={true}
          // edgesUpdatable={false}
          edgesFocusable={false}
          nodesDraggable={false}
          // nodesConnectable={false}
          nodesFocusable={false}
          draggable={true}
          panOnDrag={true}
          // elementsSelectable={false}
          // Optional if you also want to lock zooming
          zoomOnDoubleClick={false}
          onNodeClick={(_, node) => {
            nav({
              to: `/pipelines/workflows/${params.workflowId}/jobs/${node.id}`,
            });
          }}
        >
          <Controls />
          <Background />
        </ReactFlow>
      </div>
    </>
  );
}
