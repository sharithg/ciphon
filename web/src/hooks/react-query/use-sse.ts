import { useEffect } from "react";

interface WorkflowEvent {
  type: string;
}

type EventHandler = (event: WorkflowEvent) => void;

const useWorkflowEvents = (url: string, handleEvent: EventHandler) => {
  useEffect(() => {
    const evtSource = new EventSource(url);

    evtSource.onmessage = (event: MessageEvent) => {
      if (event.data) {
        const data: WorkflowEvent = JSON.parse(event.data);
        handleEvent(data);
      }
    };

    return () => {
      evtSource.close();
    };
  }, [handleEvent, url]);
};

export default useWorkflowEvents;
