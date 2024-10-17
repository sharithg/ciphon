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
import dagre from "dagre";

export const Route = createFileRoute(
  "/dashboard/_layout/pipelines/_layout/workflows/$workflowId/"
)({
  component: () => (
    <>
      <Pipeline />
    </>
  ),
});

const nodeTypes = {
  customNode: CustomNode,
};

const nodeWidth = 172;
const nodeHeight = 36;

const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  const isHorizontal = true;
  dagreGraph.setGraph({ rankdir: isHorizontal ? "LR" : "TB" });

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };
    return node;
  });

  return {
    layoutedNodes: layoutedNodes,
    layoutedEdges: edges,
  };
};

function Pipeline() {
  const params = Route.useParams();
  const { data } = useGetJobs(params.workflowId);

  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  const nav = useNavigate();

  useEffect(() => {
    if (data) {
      const sortedJobs = data.jobs.sort((a, b) => a.id.localeCompare(b.id));

      const sortedEdges = data.edges.sort((a, b) => {
        const sourceCompare = a.source.localeCompare(b.source);
        return sourceCompare !== 0
          ? sourceCompare
          : a.dest.localeCompare(b.dest);
      });

      const nodes = sortedJobs.map((v) => ({
        id: v.id,
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
        type: "customNode",
        data: { label: v.name, status: v.status },
        position: { x: 0, y: 0 },
      }));

      const edges = sortedEdges.map((v) => ({
        id: `edge-${v.dest}-${v.source}`,
        source: v.source,
        type: "smoothstep",
        target: v.dest,
        animated: false,
      }));

      const { layoutedNodes, layoutedEdges } = getLayoutedElements(
        nodes,
        edges
      );

      setNodes(layoutedNodes);
      setEdges(layoutedEdges);
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
              to: `/dashboard/pipelines/workflows/${params.workflowId}/jobs/${node.id}`,
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
