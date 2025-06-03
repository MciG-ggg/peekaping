import Layout from "@/layout";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";
import TimezoneSelector from "../../components/TimezoneSelector";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";
import { toast } from "sonner";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getSettingsKeyByKeyOptions,
  putSettingsKeyByKeyMutation,
  getSettingsKeyByKeyQueryKey,
} from "@/api/@tanstack/react-query.gen";
import React from "react";

const KeepDataPeriodSetting = () => {
  const keepDataPeriodSchema = z.object({
    value: z.coerce
      .number()
      .int()
      .min(1, { message: "Must be at least 1 day" }),
  });
  const KEEP_DATA_KEY = "KEEP_DATA_PERIOD_DAYS";
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery(
    getSettingsKeyByKeyOptions({ path: { key: KEEP_DATA_KEY } })
  );

  const form = useForm({
    resolver: zodResolver(keepDataPeriodSchema),
    defaultValues: { value: 30 },
  });

  // Reset form when data is loaded
  React.useEffect(() => {
    if (data?.data?.value) {
      form.reset({ value: Number(data.data.value) });
    }
  }, [data, form]);

  const mutation = useMutation({
    ...putSettingsKeyByKeyMutation(),
    onSuccess: () => {
      toast.success("Setting updated successfully");
      queryClient.invalidateQueries(
        getSettingsKeyByKeyQueryKey({ path: { key: KEEP_DATA_KEY } }) as any
      );
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to update setting");
    },
  });
  function onSubmit(values: { value: number }) {
    mutation.mutate({
      path: { key: KEEP_DATA_KEY },
      body: { key: KEEP_DATA_KEY, type: "int", value: String(values.value) },
    });
  }
  return (
    <Card>
      <CardHeader>
        <CardTitle>Data Retention Period</CardTitle>
        <CardDescription>
          Set how many days to keep data (KEEP_DATA_PERIOD_DAYS).
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="h-10 flex items-center">Loading...</div>
        ) : (
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="flex gap-2 items-end max-w-xs"
            >
              <FormField
                control={form.control}
                name="value"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Days</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={1}
                        {...field}
                        disabled={isLoading}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" disabled={mutation.isPending || isLoading}>
                {mutation.isPending ? "Saving..." : "Save"}
              </Button>
            </form>
          </Form>
        )}
      </CardContent>
    </Card>
  );
};

const SettingsPage = () => {
  return (
    <Layout pageName="Settings">
      <div className="max-w-xl flex flex-col gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Timezone</CardTitle>
            <CardDescription>
              Select your preferred timezone for displaying all times in the
              app.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TimezoneSelector />
          </CardContent>
        </Card>
        {/* Add more settings sections here as needed */}
        <KeepDataPeriodSetting />
      </div>
    </Layout>
  );
};

export default SettingsPage;
