import { z } from "zod";
import { TypographyH4 } from "@/components/ui/typography";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import Intervals, {
  intervalsDefaultValues,
  intervalsSchema,
} from "../shared/intervals";
import General, {
  generalDefaultValues,
  generalSchema,
} from "../shared/general";
import Notifications, {
  notificationsDefaultValues,
  notificationsSchema,
} from "../shared/notifications";
import Proxies, {
  proxiesDefaultValues,
  proxiesSchema,
} from "../shared/proxies";
import { useMonitorFormContext } from "../../context/monitor-form-context";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
import type { MonitorCreateUpdateDto, MonitorMonitorResponseDto } from "@/api";

interface PingConfig {
  host: string;
  count: number;
  packet_size: number;
}

export const pingSchema = z
  .object({
    type: z.literal("ping"),
    host: z.string().min(1, "Host is required"),
    count: z.number().min(1, "Count must be at least 1").max(10, "Count must be at most 10").optional(),
    packet_size: z.number().min(0, "Data size must be at least 0 bytes").max(65507, "Data size must be at most 65507 bytes").optional(),
  })
  .merge(generalSchema)
  .merge(intervalsSchema)
  .merge(notificationsSchema)
  .merge(proxiesSchema);

export type PingForm = z.infer<typeof pingSchema>;

export const pingDefaultValues: PingForm = {
  type: "ping",
  host: "example.com",
  count: 1,
  packet_size: 32,
  ...generalDefaultValues,
  ...intervalsDefaultValues,
  ...notificationsDefaultValues,
  ...proxiesDefaultValues,
};

export const deserialize = (data: MonitorMonitorResponseDto): PingForm => {
  let config: PingConfig = {
    host: "example.com",
    count: 1,
    packet_size: 32,
  };

  if (data.config) {
    try {
      const parsedConfig = JSON.parse(data.config);
      config = {
        host: parsedConfig.host || "example.com",
        count: parsedConfig.count || 1,
        packet_size: parsedConfig.packet_size || 32,
      };
    } catch (error) {
      console.error("Failed to parse ping monitor config:", error);
    }
  }

  return {
    type: "ping",
    name: data.name || "My Ping Monitor",
    host: config.host,
    count: config.count,
    packet_size: config.packet_size,
    interval: data.interval || 60,
    timeout: data.timeout || 16,
    max_retries: data.max_retries || 3,
    retry_interval: data.retry_interval || 60,
    resend_interval: data.resend_interval || 10,
    notification_ids: data.notification_ids || [],
    proxy_id: data.proxy_id || "",
  };
};

export const serialize = (formData: PingForm): MonitorCreateUpdateDto => {
  const config: PingConfig = {
    host: formData.host,
    count: formData.count || 1,
    packet_size: formData.packet_size || 32,
  };

  return {
    type: "ping",
    name: formData.name,
    interval: formData.interval,
    max_retries: formData.max_retries,
    retry_interval: formData.retry_interval,
    notification_ids: formData.notification_ids,
    proxy_id: formData.proxy_id,
    resend_interval: formData.resend_interval,
    timeout: formData.timeout,
    config: JSON.stringify(config),
  };
};

const PingForm = () => {
  const {
    form,
    setNotifierSheetOpen,
    setProxySheetOpen,
    isPending,
    mode,
    createMonitorMutation,
    editMonitorMutation,
    monitorId,
    monitor,
  } = useMonitorFormContext();

  const onSubmit = (data: PingForm) => {
    const payload = serialize(data);

    if (mode === "create") {
      createMonitorMutation.mutate({
        body: {
          ...payload,
          active: true,
        },
      });
    } else {
      editMonitorMutation.mutate({
        path: { id: monitorId! },
        body: {
          ...payload,
          active: monitor?.data?.active,
        },
      });
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) => onSubmit(data as PingForm))}
        className="space-y-6 max-w-[600px]"
      >
        <Card>
          <CardContent className="space-y-4">
            <General />
          </CardContent>
        </Card>

        <Card>
          <CardContent className="space-y-4">
            <TypographyH4>Ping Settings</TypographyH4>
            <FormField
              control={form.control}
              name="host"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Host</FormLabel>
                  <FormControl>
                    <Input placeholder="example.com" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="count"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Ping Count</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      placeholder="1"
                      min="1"
                      max="10"
                      {...field}
                      onChange={(e) => field.onChange(parseInt(e.target.value, 10) || 1)}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="packet_size"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Data Size (bytes)</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      placeholder="32"
                      min="0"
                      max="65507"
                      {...field}
                      onChange={(e) => field.onChange(parseInt(e.target.value, 10) || 32)}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
        </Card>

        <Card>
          <CardContent className="space-y-4">
            <Notifications onNewNotifier={() => setNotifierSheetOpen(true)} />
          </CardContent>
        </Card>

        <Card>
          <CardContent className="space-y-4">
            <Proxies onNewProxy={() => setProxySheetOpen(true)} />
          </CardContent>
        </Card>

        <Card>
          <CardContent className="space-y-4">
            <Intervals />
          </CardContent>
        </Card>

        <Button type="submit">
          {isPending && <Loader2 className="animate-spin" />}
          {mode === "create" ? "Create" : "Update"}
        </Button>
      </form>
    </Form>
  );
};

export default PingForm;
