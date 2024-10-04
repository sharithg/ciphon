import { Handle, Position } from "@xyflow/react";

const CustomNode = ({
  data,
}: {
  data: { label: string; description: string };
}) => {
  return (
    <div className="p-3 border border-black-300 rounded-md bg-black-100 text-center shadow-md">
      <strong>{data.label}</strong>
      <div className="mt-1">
        <small>{data.description}</small>
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
