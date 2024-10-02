import { useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

const ConnectionStatus = {
  [ReadyState.CONNECTING]: "Connecting",
  [ReadyState.OPEN]: "Open",
  [ReadyState.CLOSING]: "Closing",
  [ReadyState.CLOSED]: "Closed",
  [ReadyState.UNINSTANTIATED]: "Uninstantiated",
} as const;

export const useWebsocket = () => {
  const [socketUrl] = useState("ws://localhost:8000/ws");

  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

  const connectionStatus = ConnectionStatus[readyState];

  return {
    connectionStatus,
    sendMessage,
    lastMessage,
  };
};
