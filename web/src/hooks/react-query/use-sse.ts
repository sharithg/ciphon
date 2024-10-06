import { useEffect } from "react";
import { API_URL } from "./constants";

interface WorkflowEvent {
  type: string;
}

type EventHandler = (event: WorkflowEvent) => void;

const useWorkflowEvents = (handleEvent: EventHandler) => {
  useEffect(() => {
    const evtSource = new EventSource(`${API_URL}/workflows/run-events`);

    evtSource.onmessage = (event: MessageEvent) => {
      if (event.data) {
        const data: WorkflowEvent = JSON.parse(event.data);
        handleEvent(data);
      }
    };

    return () => {
      evtSource.close();
    };
  }, [handleEvent]);
};

export default useWorkflowEvents;
