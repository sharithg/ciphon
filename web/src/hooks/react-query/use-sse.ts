import { useEffect } from "react";

export const useSse = (input: {
  url: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onMessage: (data: any) => void;
}) => {
  useEffect(() => {
    const evtSource = new EventSource(input.url);
    evtSource.onmessage = (event) => {
      if (event.data) {
        input.onMessage(JSON.parse(event.data));
      }
    };

    return () => {
      evtSource.close();
    };
  }, [input]);
};
