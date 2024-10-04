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

export const Route = createLazyFileRoute("/pipelines/$pipelineId")({
  component: Pipeline,
});

const nodeTypes = {
  customNode: CustomNode,
};

const nodes: Node[] = [
  {
    id: "horizontal-1",
    sourcePosition: Position.Right,
    type: "customNode",
    data: { label: "Input" },
    position: { x: 0, y: 80 },
  },
  {
    id: "horizontal-2",
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: { label: "A Node" },
    position: { x: 250, y: 0 },
    type: "customNode",
  },
  {
    id: "horizontal-3",
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: { label: "Node 3" },
    position: { x: 250, y: 160 },
    type: "customNode",
  },
  {
    id: "horizontal-4",
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: { label: "Node 4" },
    position: { x: 500, y: 0 },
    type: "customNode",
  },
  {
    id: "horizontal-5",
    sourcePosition: Position.Top,
    targetPosition: Position.Bottom,
    data: { label: "Node 5" },
    position: { x: 500, y: 100 },
    type: "customNode",
  },
  {
    id: "horizontal-6",
    sourcePosition: Position.Bottom,
    targetPosition: Position.Top,
    data: { label: "Node 6" },
    position: { x: 500, y: 230 },
    type: "customNode",
  },
  {
    id: "horizontal-7",
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: { label: "Node 7" },
    position: { x: 750, y: 50 },
    type: "customNode",
  },
  {
    id: "horizontal-8",
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: { label: "Node 8" },
    position: { x: 750, y: 300 },
    type: "customNode",
  },
];

const edges: Edge[] = [
  {
    id: "horizontal-e1-2",
    source: "horizontal-1",
    type: "smoothstep",
    target: "horizontal-2",
  },
  {
    id: "horizontal-e1-3",
    source: "horizontal-1",
    type: "smoothstep",
    target: "horizontal-3",
  },
  {
    id: "horizontal-e1-4",
    source: "horizontal-2",
    type: "smoothstep",
    target: "horizontal-4",
  },
  {
    id: "horizontal-e3-5",
    source: "horizontal-3",
    type: "smoothstep",
    target: "horizontal-5",
  },
  {
    id: "horizontal-e3-6",
    source: "horizontal-3",
    type: "smoothstep",
    target: "horizontal-6",
  },
  {
    id: "horizontal-e5-7",
    source: "horizontal-5",
    type: "smoothstep",
    target: "horizontal-7",
  },
  {
    id: "horizontal-e6-8",
    source: "horizontal-6",
    type: "smoothstep",
    target: "horizontal-8",
  },
];

function Pipeline() {
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
