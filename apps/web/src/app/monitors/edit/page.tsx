import {
  getMonitorsByIdOptions,
  getMonitorsByIdQueryKey,
  getMonitorsInfiniteQueryKey,
  getNotificationsQueryKey,
  getProxiesQueryKey,
  putMonitorsByIdMutation,
} from "@/api/@tanstack/react-query.gen";
import Layout from "@/layout";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";
import CreateEditFields from "../components/create-edit-fields";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { Form } from "@/components/ui/form";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import CreateEditNotifier from "@/app/notifiers/create-edit-notifier";
import CreateEditProxy from "@/app/proxies/components/create-edit-proxy";
import { generalDefaultValues, generalSchema } from "../components/general";
import { httpDefaultValues, httpSchema } from "../components/http";
import { intervalsDefaultValues } from "../components/intervals";

const formSchema = z.object({
  general: generalSchema,

  intervals: z.object({
    interval: z.coerce
      .number()
      .min(1, { message: "Interval must be at least 1 second" })
      .max(3600, { message: "Interval must be less than 1 hour" }),
    max_retries: z.coerce.number(),
    retry_interval: z.coerce.number(),
    timeout: z.coerce.number(),
    resend_interval: z.coerce.number(),
  }),

  notifications: z.object({
    notification_ids: z.array(z.string()),
  }),

  proxies: z.object({
    proxy_id: z.string().optional(),
  }),

  http: httpSchema,
});

export type Form = z.infer<typeof formSchema>;

const formDefaultValues: Form = {
  general: generalDefaultValues,

  intervals: intervalsDefaultValues,

  notifications: {
    notification_ids: [],
  },

  proxies: {
    proxy_id: undefined,
  },

  // type specific
  http: httpDefaultValues,
};

const EditMonitor = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [notifierSheetOpen, setNotifierSheetOpen] = useState(false);
  const [proxySheetOpen, setProxySheetOpen] = useState(false);

  const { data: monitor } = useQuery({
    ...getMonitorsByIdOptions({
      path: {
        id: id!,
      },
    }),
    enabled: !!id,
  });

  const form = useForm<Form>({
    defaultValues: formDefaultValues,
    resolver: zodResolver(formSchema),
  });

  useEffect(() => {
    if (monitor?.data) {
      const monitorType = monitor.data.type as "http"; // Cast to expected type
      const config = monitor.data.config ? JSON.parse(monitor.data.config) : {};
      form.reset({
        general: {
          name: monitor.data.name || "",
          type: monitorType,
        },
        intervals: {
          interval: monitor.data.interval || 60,
          max_retries: monitor.data.max_retries || 3,
          retry_interval: monitor.data.retry_interval || 60,
          resend_interval: monitor.data.resend_interval || 10,
          timeout: monitor.data.timeout || 16,
        },
        notifications: {
          notification_ids: monitor.data.notification_ids || [],
        },
        proxies: {
          proxy_id: monitor.data.proxy_id || undefined,
        },
        [monitorType]: config,
      });
    }
  }, [form, monitor]);

  const updateMonitorMutation = useMutation({
    ...putMonitorsByIdMutation(),
    onSuccess: () => {
      toast.success("Monitor updated successfully");
      queryClient.invalidateQueries({
        queryKey: getMonitorsInfiniteQueryKey(),
      });
      queryClient.invalidateQueries({
        queryKey: getMonitorsByIdQueryKey({ path: { id: id! } }),
      });
      navigate(`/monitors/${id}`, { replace: true });
    },
    onError: (error) => {
      console.log(error);
      if (isAxiosError(error)) {
        toast.error(error.response?.data.message || error.message);
      } else {
        console.log(error);
      }
    },
  });

  const onSubmit = (formData: Form) => {
    const { general, intervals, notifications, proxies, ...rest } = formData;
    const typeSpecific = rest[general.type] || {};
    updateMonitorMutation.mutate({
      path: {
        id: id!,
      },
      body: {
        ...general,
        ...intervals,
        ...notifications,
        proxy_id: proxies.proxy_id || undefined,
        config: JSON.stringify(typeSpecific),
        active: monitor?.data?.active,
      },
    });
  };

  return (
    <Layout pageName={`Edit Monitor | ${monitor?.data?.name}`}>
      <div className="flex flex-col gap-4">
        <p className="text-gray-500">
          Create a new monitor to start tracking your website's performance.
        </p>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-6 max-w-[600px]"
          >
            <CreateEditFields
              onNewNotifier={() => setNotifierSheetOpen(true)}
              onNewProxy={() => setProxySheetOpen(true)}
            />

            <Button type="submit">
              {updateMonitorMutation.isPending && (
                <Loader2 className="animate-spin" />
              )}
              Update
            </Button>
          </form>
        </Form>
      </div>

      <Sheet open={notifierSheetOpen} onOpenChange={setNotifierSheetOpen}>
        <SheetContent
          className="p-4 overflow-y-auto"
          onInteractOutside={(event) => event.preventDefault()}
        >
          <CreateEditNotifier
            onSuccess={async (newNotifier) => {
              setNotifierSheetOpen(false);
              queryClient.invalidateQueries({
                queryKey: getNotificationsQueryKey(),
              });
              form.setValue("notifications.notification_ids", [
                ...(form.getValues("notifications.notification_ids") || []),
                newNotifier.id,
              ]);
            }}
          />
        </SheetContent>
      </Sheet>

      <Sheet open={proxySheetOpen} onOpenChange={setProxySheetOpen}>
        <SheetContent
          className="p-4 overflow-y-auto"
          onInteractOutside={(event) => event.preventDefault()}
        >
          <CreateEditProxy
            onSuccess={() => {
              setProxySheetOpen(false);
              queryClient.invalidateQueries({
                queryKey: getProxiesQueryKey(),
              });
              // The new proxy will appear in the dropdown after refetching
            }}
          />
        </SheetContent>
      </Sheet>
    </Layout>
  );
};

export default EditMonitor;
