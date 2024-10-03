import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { z } from "zod";
import { toast } from "../use-toast";
import { Node } from "@/@types/api";
import { API_URL } from "./constants";

export const nodeSchema = z.object({
  name: z.string().min(2, {
    message: "Name must be at least 2 characters long.",
  }),
  user: z.string().min(2, {
    message: "User name must be at least 2 characters long.",
  }),
  host: z.string().min(2, {
    message: "Host must be at least 2 characters long.",
  }),
  port: z.coerce.number(),
  pem: z.instanceof(File, {
    message: "Pem file is required and must be a valid file.",
  }),
});

export const useAddNewNode = (input: { onSuccess?: () => Promise<void> }) => {
  const mutation = useMutation({
    mutationFn: (newNode: z.infer<typeof nodeSchema>) => {
      const formData = new FormData();
      formData.append("pem", newNode.pem);
      formData.append("host", newNode.host);
      formData.append("name", newNode.name);
      formData.append("user", newNode.user);
      formData.append("port", newNode.port.toString());

      return axios.post(`${API_URL}/node`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
    },
    onSuccess: () => {
      toast({
        title: "Succesfully added new node",
      });

      if (input.onSuccess) {
        input.onSuccess();
      }
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
    queryKey: ["nodes"],
    queryFn: () => axios.get<Node[]>(`${API_URL}/nodes`),
    refetchInterval: 3000,
  });
};
