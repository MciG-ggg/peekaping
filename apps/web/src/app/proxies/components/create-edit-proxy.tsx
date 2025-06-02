import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { z } from "zod";
import { useForm, useWatch } from "react-hook-form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { postProxiesMutation, getProxiesInfiniteQueryKey } from "@/api/@tanstack/react-query.gen";
import { Loader2 } from "lucide-react";
import type { ProxyCreateUpdateDto } from "@/api/types.gen";

const formSchema = z.object({
  protocol: z.enum(["http", "https", "socks", "socks5", "socks5h", "socks4"]),
  host: z.string().min(1, { message: "Host is required" }),
  port: z
    .number()
    .min(1, { message: "Port must be at least 1" })
    .max(65535, { message: "Port must be between 1 and 65535" }),
  auth: z.boolean(),
  username: z.string().optional(),
  password: z.string().optional(),
});

type Form = z.infer<typeof formSchema>;
export type { Form };

const proxyProtocols = [
  { type: "http", description: "HTTP" },
  { type: "https", description: "HTTPS" },
  { type: "socks", description: "SOCKS" },
  { type: "socks5", description: "SOCKS v5" },
  { type: "socks5h", description: "SOCKS v5 (+DNS)" },
  { type: "socks4", description: "SOCKS v4" },
];

type CreateEditProxyProps = {
  onSuccess?: () => void;
  initialValues?: Form;
  isEdit?: boolean;
  isLoading?: boolean;
  onSubmit?: (data: Form) => void;
};

export default function CreateEditProxy({
  onSuccess,
  initialValues,
  isEdit = false,
  isLoading = false,
  onSubmit: externalSubmit
}: CreateEditProxyProps) {
  const queryClient = useQueryClient();

  const form = useForm<Form>({
    defaultValues: initialValues || {
      protocol: "https",
      host: "",
      port: 1,
      auth: false,
      username: "",
      password: "",
    },
    resolver: zodResolver(formSchema),
  });

  const { isSubmitting } = form.formState;

  const mutation = useMutation({
    ...postProxiesMutation(),
    onSuccess: () => {
      toast.success(`Proxy ${isEdit ? 'updated' : 'created'} successfully`);
      queryClient.invalidateQueries({ queryKey: getProxiesInfiniteQueryKey() });
      if (onSuccess) onSuccess();
    },
    onError: (error) => {
      console.error(`Error ${isEdit ? 'updating' : 'creating'} proxy:`, error);
      toast.error(error.message || `Failed to ${isEdit ? 'update' : 'create'} proxy`);
    },
  });

  const onSubmit = (data: Form) => {
    // If there's an external submit handler (for editing), use that
    if (externalSubmit) {
      externalSubmit(data);
      return;
    }

    // Otherwise handle creation
    const proxyData: ProxyCreateUpdateDto = {
      protocol: data.protocol,
      host: data.host,
      port: data.port,
      auth: data.auth,
      username: data.auth ? data.username : undefined,
      password: data.auth ? data.password : undefined,
    };

    // Make API call
    mutation.mutate({
      body: proxyData,
    });
  };

  const auth = useWatch({
    control: form.control,
    name: "auth",
  });

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="space-y-6 max-w-[600px]"
      >
        <FormField
          control={form.control}
          name="protocol"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Proxy Protocol</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select proxy protocol" />
                  </SelectTrigger>
                </FormControl>

                <SelectContent>
                  {proxyProtocols.map((item) => (
                    <SelectItem key={item.type} value={item.type}>
                      {item.description}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormItem>
          <FormLabel>Proxy Server</FormLabel>
          <div className="flex space-x-4">
            <FormField
              control={form.control}
              name="host"
              render={({ field }) => (
                <>
                  <Input placeholder="Server address" {...field} />
                  <FormMessage />
                </>
              )}
            />

            <FormField
              control={form.control}
              name="port"
              render={({ field }) => (
                <>
                  <Input
                    placeholder="Port"
                    {...field}
                    type="number"
                    value={field.value}
                    onChange={e => field.onChange(Number(e.target.value))}
                  />
                  <FormMessage />
                </>
              )}
            />
          </div>
        </FormItem>

        <FormField
          control={form.control}
          name="auth"
          render={({ field }) => (
            <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
              <div className="space-y-0.5">
                <FormLabel>Proxy server has authentication</FormLabel>
              </div>

              <FormControl>
                <Switch
                  checked={field.value}
                  onCheckedChange={field.onChange}
                  aria-readonly
                />
              </FormControl>
            </FormItem>
          )}
        />

        {auth && (
          <>
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>User</FormLabel>
                  <FormControl>
                    <Input placeholder="User" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input placeholder="Password" {...field} type="password" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </>
        )}

        <Button
          type="submit"
          disabled={isSubmitting || mutation.isPending || isLoading}
        >
          {(isSubmitting || mutation.isPending || isLoading) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          {(isSubmitting || mutation.isPending || isLoading) ? "Saving..." : isEdit ? "Update" : "Save"}
        </Button>
      </form>
    </Form>
  );
}
