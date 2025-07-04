import { Input } from "@/components/ui/input";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { z } from "zod";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { useFormContext } from "react-hook-form";

export const schema = z.object({
  type: z.literal("teams"),
  webhook_url: z.string().url({ message: "Valid webhook URL is required" }),
  use_template: z.boolean().optional(),
  template: z.string().optional(),
});

export type TeamsFormValues = z.infer<typeof schema>;

export const defaultValues: TeamsFormValues = {
  type: "teams",
  webhook_url: "",
  use_template: false,
  template: "",
};

export const displayName = "Microsoft Teams";

export default function TeamsForm() {
  const form = useFormContext();
  const useTemplate = form.watch("use_template");

  const templatePlaceholder = `Example template:
🔔 Monitor Alert: {{ monitor.name }}
Status: {{ status }}
Message: {{ msg }}
Time: {{ heartbeat.time }}`;

  return (
    <>
      <FormField
        control={form.control}
        name="webhook_url"
        render={({ field }) => (
          <FormItem>
            <FormLabel>
              Webhook URL <span className="text-red-500">*</span>
            </FormLabel>
            <FormControl>
              <Input
                placeholder="https://your-teams-webhook-url"
                type="url"
                required
                {...field}
              />
            </FormControl>
            <FormDescription>
              <span className="text-red-500">*</span> Required
              <br />
              <span className="mt-2 block">
                Learn how to get a Teams webhook URL:{" "}
                <a
                  href="https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="underline text-blue-600"
                >
                  Microsoft Teams Documentation
                </a>
              </span>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="use_template"
        render={({ field }) => (
          <FormItem>
            <div className="flex items-center gap-2">
              <FormControl>
                <Switch
                  checked={field.value || false}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <FormLabel>Use Custom Template</FormLabel>
            </div>
            <FormDescription>
              Enable to use a custom message template instead of the default format.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {useTemplate && (
        <FormField
          control={form.control}
          name="template"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Custom Template</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={templatePlaceholder}
                  className="min-h-[150px] font-mono text-sm"
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Customize the notification message format. Available variables:
                <br />
                <code className="text-pink-500 ml-1">{"{{ monitor.name }}"}</code> - Monitor name,{" "}
                <code className="text-pink-500">{"{{ status }}"}</code> - Status (UP/DOWN),{" "}
                <code className="text-pink-500">{"{{ msg }}"}</code> - Status message
                <br />
                <code className="text-pink-500">{"{{ heartbeat.time }}"}</code> - Timestamp,{" "}
                <code className="text-pink-500">{"{{ monitor.type }}"}</code> - Monitor type
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      )}
    </>
  );
}
