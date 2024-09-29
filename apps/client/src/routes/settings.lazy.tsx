import { createLazyFileRoute } from "@tanstack/react-router";
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
  NodeWithStatus,
  useAddNewNode,
  useGetNodes,
} from "../hooks/use-nodes";

export const Route = createLazyFileRoute("/settings")({
  component: Settings,
});

function Settings() {
  const mutation = useAddNewNode();

  const nodes = useGetNodes();

  const nodesWithStatus: NodeWithStatus[] = (nodes.data?.data ?? []).map(
    (n) => ({
      ...n,
      status: "Healthy",
    })
  );

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
                                node.status === "Healthy"
                                  ? "bg-green-100 text-green-800"
                                  : "bg-red-100 text-red-800"
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
                    <FormField
                      control={form.control}
                      name="name"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Name</FormLabel>
                          <FormControl>
                            <Input placeholder="name" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="host"
                      render={({ field }) => (
                        <FormItem>
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
                      name="user"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>User</FormLabel>
                          <FormControl>
                            <Input placeholder="user" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
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
