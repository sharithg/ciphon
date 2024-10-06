import { createLazyFileRoute, Link } from "@tanstack/react-router";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
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
import { useGetJobs } from "../../../hooks/react-query/use-workflows";
import { useEffect, useState } from "react";
import useWorkflowEvents from "../../../hooks/react-query/use-sse";

export const Route = createLazyFileRoute("/pipelines/workflows/$workflowId")({
  component: Pipeline,
});

const nodeTypes = {
  customNode: CustomNode,
};

const edges: Edge[] = [];

function Pipeline() {
  const params = Route.useParams();
  const { data, refetch } = useGetJobs(params.workflowId);

  const [nodes, setNodes] = useState<Node[]>([]);

  useWorkflowEvents((event) => {
    if (event.type === "job") {
      refetch();
    }
  });

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

  console.log({ data: data });

  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/">Pipelines</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Breadcrumb</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
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
            console.log(node.id);
          }}
        >
          <Controls />
          <Background />
        </ReactFlow>
      </div>
    </>
  );
}
