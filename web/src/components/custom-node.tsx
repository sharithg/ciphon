import { Handle, Position } from "@xyflow/react";
import { CheckCircle, XCircle, Loader2 } from "lucide-react";
import { cn } from "../lib/utils";

type TStatus = "running" | "success" | "failed";

function getStatusIcon(status: TStatus) {
  switch (status) {
    case "running":
      return <Loader2 className="animate-spin" />;
    case "success":
      return <CheckCircle className="text-green-500" />;
    case "failed":
      return <XCircle className="text-red-500" />;
    default:
      return null;
  }
}

function getColor(status: TStatus) {
  switch (status) {
    case "running":
      return "bg-blue-400";
    case "success":
      return "bg-green-400";
    case "failed":
      return "bg-red-400";
    default:
      return "bg-black-400";
  }
}

const CustomNode = ({
  data,
}: {
  data: {
    label: string;
    description: string;
    status: TStatus;
  };
}) => {
  const StatusIcon = getStatusIcon(data.status);
  const color = getColor(data.status);

  return (
    <div
      className={cn(
        "p-3 border border-black-300 rounded-md text-center shadow-md",
        color
      )}
    >
      <div className="flex items-center justify-center gap-2 mb-1">
        {StatusIcon}
        <strong className="text-sm">{data.label}</strong>
      </div>
      <div className="mt-1">
        <small className="text-xs">{data.description}</small>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        id="a"
        className="bg-gray-600"
      />
      <Handle
        type="target"
        position={Position.Left}
        id="b"
        className="bg-gray-600"
      />
    </div>
  );
};

export default CustomNode;
