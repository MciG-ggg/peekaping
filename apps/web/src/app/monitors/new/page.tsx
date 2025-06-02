import {
  getMonitorsInfiniteQueryKey,
  postMonitorsMutation,
  getNotificationsQueryKey,
  getProxiesQueryKey,
} from "@/api/@tanstack/react-query.gen";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { isAxiosError } from "axios";
import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import CreateEditFields from "../components/create-edit-fields";
import { Sheet, SheetContent } from "@/components/ui/sheet";
import { useState } from "react";
import CreateEditNotifier from "@/app/notifiers/create-edit-notifier";
import CreateEditProxy from "@/app/proxies/components/create-edit-proxy";
import { httpDefaultValues, httpSchema } from "../components/http";
import { generalDefaultValues, generalSchema } from "../components/general";
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

type Form = z.infer<typeof formSchema>;

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

const NewMonitor = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [notifierSheetOpen, setNotifierSheetOpen] = useState(false);
  const [proxySheetOpen, setProxySheetOpen] = useState(false);

  const form = useForm<Form>({
    defaultValues: formDefaultValues,
    resolver: zodResolver(formSchema),
  });

  const createMonitorMutation = useMutation({
    ...postMonitorsMutation(),
    onSuccess: () => {
      toast.success("Monitor created successfully");
      queryClient.invalidateQueries({
        queryKey: getMonitorsInfiniteQueryKey(),
      });
      navigate("/monitors");
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

  const onSubmit = (data: Form) => {
    const { general, intervals, notifications, proxies, ...rest } = data;
    const typeSpecific = rest[general.type] || {};

    createMonitorMutation.mutate({
      body: {
        ...general,
        ...intervals,
        ...notifications,
        proxy_id: proxies.proxy_id || undefined,
        config: JSON.stringify(typeSpecific),
      },
    });
  };

  return (
    <Layout pageName="New Monitor">
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
              {createMonitorMutation.isPending && (
                <Loader2 className="animate-spin" />
              )}
              Save
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

export default NewMonitor;
