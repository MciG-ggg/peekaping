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
import { useForm } from "react-hook-form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { CardTitle } from "@/components/ui/card";
import * as SmtpForm from "./integrations/smtp-form";
import { Button } from "@/components/ui/button";
import {
  getNotificationsInfiniteQueryKey,
  postNotificationsMutation,
} from "@/api/@tanstack/react-query.gen";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import * as TelegramForm from "./integrations/telegram-form";
import * as WebhookForm from "./integrations/webhook-form";
import * as SlackForm from "./integrations/slack-form";
import { useEffect } from "react";

const typeFormRegistry = {
  smtp: SmtpForm,
  telegram: TelegramForm,
  webhook: WebhookForm,
  slack: SlackForm,
};

const baseSchema = z
  .object({
    name: z.string().min(1, {
      message: "Name is required",
    }),
  })
  .and(
    z.discriminatedUnion("type", [
      SmtpForm.schema,
      TelegramForm.schema,
      WebhookForm.schema,
      SlackForm.schema,
    ])
  );

type BaseFormValues = z.infer<typeof baseSchema>;

// validate map components
Object.values(typeFormRegistry).forEach((component) => {
  if (typeof component.default !== "function") {
    throw new Error("Type components must be exported as default");
  }
  if (!component.displayName) {
    throw new Error("Type components must have a displayName");
  }
  if (!component.schema) {
    throw new Error("Type components must have a schema");
  }
  if (!component.defaultValues) {
    throw new Error("Type components must have default values");
  }
});

const notificationTypes = Object.keys(typeFormRegistry).map((key) => ({
  label: typeFormRegistry[key as keyof typeof typeFormRegistry].displayName,
  value: key,
}));

export default function CreateEditNotifier({
  onSuccess,
  onSubmit,
  initialValues = {
    name: "My notifier",
    ...SmtpForm.defaultValues,
  },
  isLoading = false,
  isEdit = false,
}: {
  onSuccess?: (notifier: { id: string }) => void;
  onSubmit?: (data: { name: string; type: string; config: string }) => void;
  initialValues?: BaseFormValues;
  isLoading?: boolean;
  isEdit?: boolean;
} = {}) {
  const queryClient = useQueryClient();
  // Form for name and type
  const baseForm = useForm<BaseFormValues>({
    resolver: zodResolver(baseSchema),
    defaultValues: initialValues,
    // shouldUnregister: true,
  });

  console.log(baseForm.formState.errors);

  const type = baseForm.watch("type");

  const TypeFormComponent = typeFormRegistry[type]?.default || null;

  useEffect(() => {
    if (type === initialValues?.type) return;
    if (!type) return;

    const values = baseForm.getValues();
    baseForm.reset({
      ...values,
      ...(typeFormRegistry[type].defaultValues || {}),
    });
  }, [type]);

  const createNotifierMutation = useMutation({
    ...postNotificationsMutation(),
    onSuccess: (data) => {
      toast.success("Notifier created successfully");
      if (onSuccess && data?.data?.id) {
        onSuccess({ id: data.data.id });
      }
      queryClient.invalidateQueries({
        queryKey: getNotificationsInfiniteQueryKey(),
      });
    },
    onError: () => {
      toast.error("Failed to create notifier");
    },
  });

  // Handle submit
  const handleSubmit = (data: BaseFormValues) => {
    if (onSubmit) {
      // For edit mode, use the custom onSubmit handler
      const { name, type, ...typeConfig } = data;
      onSubmit({
        name,
        type,
        config: JSON.stringify(typeConfig),
      });
    } else {
      // For create mode, use the default mutation
      const { name, type, ...typeConfig } = data;
      createNotifierMutation.mutate({
        body: {
          name,
          type,
          config: JSON.stringify(typeConfig),
          active: true,
          is_default: false,
        },
      });
    }
  };

  return (
    <div className="flex flex-col gap-6 max-w-[600px]">
      <CardTitle className="text-xl">
        {isEdit ? "Edit" : "Create"} Notifier
      </CardTitle>

      <Form {...baseForm}>
        <form
          onSubmit={baseForm.handleSubmit(handleSubmit)}
          className="space-y-6 max-w-[600px]"
        >
          <FormItem>
            <FormLabel>Notifier Type</FormLabel>
            <Select
              onValueChange={(val) => {
                if (!val) return;
                baseForm.setValue("type", val as "smtp" | "telegram");
              }}
              value={type}
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select notifier type" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {notificationTypes.map((item) => (
                  <SelectItem key={item.value} value={item.value}>
                    {item.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>

          <FormField
            control={baseForm.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Friendly name</FormLabel>
                <FormControl>
                  <Input placeholder="Friendly name" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          {TypeFormComponent && <TypeFormComponent />}

          <Button
            type="submit"
            disabled={isLoading || createNotifierMutation.isPending}
          >
            {isLoading || createNotifierMutation.isPending
              ? "Saving..."
              : "Save"}
          </Button>
        </form>
      </Form>
    </div>
  );
}
