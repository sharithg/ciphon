import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Loader2, PlusCircle } from "lucide-react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  nodeSchema,
  useAddNewNode,
  useGetNodes,
} from "@/hooks/react-query/use-nodes";
import { useWebsocket } from "@/hooks/use-websocket";
import { useEffect } from "react";

export const Route = createFileRoute("/dashboard/_layout/settings")({
  component: Settings,
});

// ? "bg-green-100 text-green-800"
//                               : "bg-red-100 text-red-800"

const StatusToColor = {
  provisioning: "bg-orange-100 text-orange-800",
  error: "bg-red-100 text-red-800",
  healthy: "bg-green-100 text-green-800",
} as const;

function Settings() {
  const nodes = useGetNodes();
  const mutation = useAddNewNode({
    onSuccess: async () => {
      await nodes.refetch();
    },
  });
  const { connectionStatus, sendMessage } = useWebsocket();

  console.log({ connectionStatus });

  useEffect(() => {
    if (connectionStatus === "Open") {
      console.log("sending message");
      sendMessage("asdsd");
    }
  }, [sendMessage, connectionStatus]);

  const nodesWithStatus = nodes.data ?? [];

  function onSubmit(values: z.infer<typeof nodeSchema>) {
    mutation.mutate(values);
  }

  const form = useForm<z.infer<typeof nodeSchema>>({
    resolver: zodResolver(nodeSchema),
    defaultValues: {},
  });

  return (
    <div className="container mx-auto p-2">
      <Tabs defaultValue="nodes" className="space-y-4">
        <TabsList className="mb-4">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="nodes">Nodes</TabsTrigger>
          <TabsTrigger value="security">Security</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
        </TabsList>
        <TabsContent value="nodes" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <Card className="md:col-span-2">
              <CardHeader>
                <CardTitle>Manage Nodes</CardTitle>
                <CardDescription>
                  View and manage your registered nodes (VMs).
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ScrollArea className="h-[400px] w-full rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-[200px]">Name</TableHead>
                        <TableHead>Hostname</TableHead>
                        <TableHead className="text-right">Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {nodesWithStatus.map((node) => (
                        <TableRow key={node.id}>
                          <TableCell className="font-medium">
                            {node.name}
                          </TableCell>
                          <TableCell>{node.host}</TableCell>
                          <TableCell className="text-right">
                            <span
                              className={`inline-block px-2 py-1 rounded-full text-xs font-semibold ${
                                StatusToColor[
                                  node.status as keyof typeof StatusToColor
                                ]
                              }`}
                            >
                              {node.status}
                            </span>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </ScrollArea>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Register New Node</CardTitle>
                <CardDescription>
                  Add a new node to your dashboard.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Form {...form}>
                  <form
                    onSubmit={form.handleSubmit(onSubmit)}
                    className="space-y-8"
                  >
                    <div className="flex space-x-4">
                      <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem className="flex-1">
                            <FormLabel>Nickname</FormLabel>
                            <FormControl>
                              <Input placeholder="name" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="user"
                        render={({ field }) => (
                          <FormItem className="flex-1">
                            <FormLabel>Username</FormLabel>
                            <FormControl>
                              <Input placeholder="user" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                    <div className="flex space-x-4">
                      <FormField
                        control={form.control}
                        name="host"
                        render={({ field }) => (
                          <FormItem className="flex-1">
                            <FormLabel>Host</FormLabel>
                            <FormControl>
                              <Input placeholder="host" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={form.control}
                        name="port"
                        render={({ field }) => (
                          <FormItem className="w-1/4">
                            <FormLabel>Port</FormLabel>
                            <FormControl>
                              <Input
                                placeholder="port"
                                type="number"
                                min="1"
                                max="65535"
                                {...field}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <FormField
                      control={form.control}
                      name="pem"
                      render={({
                        // eslint-disable-next-line @typescript-eslint/no-unused-vars
                        field: { value, onChange, ...fieldProps },
                      }) => (
                        <FormItem>
                          <FormLabel>Pem file</FormLabel>
                          <FormControl>
                            <Input
                              {...fieldProps}
                              placeholder="Pem file"
                              type="file"
                              accept=".pem,application/x-pem-file,application/x-x509-ca-cert,application/octet-stream"
                              onChange={async (event) => {
                                const content =
                                  event.target.files?.length &&
                                  event.target.files[0];
                                onChange(content);
                              }}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <Button type="submit" disabled={mutation.isLoading}>
                      {mutation.isLoading ? (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      ) : (
                        <PlusCircle className="mr-2 h-4 w-4" />
                      )}
                      Register Node
                    </Button>
                  </form>
                </Form>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
