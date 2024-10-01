import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { z } from "zod";
import { toast } from "./use-toast";
import { Node } from "@/@types/api";

export const nodeSchema = z.object({
  name: z.string().min(2, {
    message: "name must be at least 2 characters.",
  }),
  user: z.string().min(2, {
    message: "name must be at least 2 characters.",
  }),
  host: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  pem: z.custom<File>((v) => v instanceof File, {
    message: "Pem is required",
  }),
});

export type NodeWithStatus = Node & { status: "Healthy" | "Unhealthy" };

export const useAddNewNode = () => {
  const mutation = useMutation({
    mutationFn: (newNode: z.infer<typeof nodeSchema>) => {
      const formData = new FormData();
      formData.append("pem", newNode.pem);
      formData.append("host", newNode.host);
      formData.append("name", newNode.name);
      formData.append("user", newNode.user);

      return axios.post("http://localhost:8000/node", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
    },
    onSuccess: () => {
      toast({
        title: "Succesfully added new node",
      });
    },
    onError: () => {
      toast({
        title: "Error adding new node",
        variant: "destructive",
      });
    },
  });
  return mutation;
};

export const useGetNodes = () => {
  return useQuery({
    queryKey: ["todos"],
    queryFn: () => axios.get<Node[]>("http://localhost:8000/nodes"),
  });
};
