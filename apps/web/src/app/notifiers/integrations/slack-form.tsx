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
import { useFormContext } from "react-hook-form";

export const schema = z.object({
  type: z.literal("slack"),
  slack_webhook_url: z.string().url({ message: "Valid webhook URL is required" }),
  slack_username: z.string().optional(),
  slack_icon_emoji: z.string().optional(),
  slack_channel: z.string().optional(),
  slack_rich_message: z.boolean().optional(),
  slack_channel_notify: z.boolean().optional(),
});

export type SlackFormValues = z.infer<typeof schema>;

export const defaultValues: SlackFormValues = {
  type: "slack",
  slack_webhook_url: "",
  slack_username: "",
  slack_icon_emoji: "",
  slack_channel: "",
  slack_rich_message: false,
  slack_channel_notify: false,
};

export const displayName = "Slack";

export default function SlackForm() {
  const form = useFormContext();

  return (
    <>
      <FormField
        control={form.control}
        name="slack_webhook_url"
        render={({ field }) => (
          <FormItem>
            <FormLabel>
              Webhook URL <span className="text-red-500">*</span>
            </FormLabel>
            <FormControl>
              <Input
                placeholder="https://hooks.slack.com/services/..."
                type="url"
                required
                {...field}
              />
            </FormControl>
            <FormDescription>
              <span className="text-red-500">*</span> Required
              <br />
              <span className="mt-2 block">
                Learn more about webhooks:{" "}
                <a
                  href="https://api.slack.com/messaging/webhooks"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="underline text-blue-600"
                >
                  https://api.slack.com/messaging/webhooks
                </a>
              </span>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slack_username"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input placeholder="Uptime Monitor" {...field} />
            </FormControl>
            <FormDescription>
              The username that will appear as the sender of the message. If not specified, the default bot name will be used.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slack_icon_emoji"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Icon Emoji</FormLabel>
            <FormControl>
              <Input placeholder=":warning:" {...field} />
            </FormControl>
            <FormDescription>
              Emoji to use as the icon for this message. Must be in the format :emoji_name:
              <br />
              <span className="mt-2 block">
                Emoji cheat sheet:{" "}
                <a
                  href="https://www.webfx.com/tools/emoji-cheat-sheet/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="underline text-blue-600"
                >
                  https://www.webfx.com/tools/emoji-cheat-sheet/
                </a>
              </span>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slack_channel"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Channel Name</FormLabel>
            <FormControl>
              <Input placeholder="#general" {...field} />
            </FormControl>
            <FormDescription>
              The channel where the message will be sent. If not specified, the message will be sent to the default channel configured in the webhook.
              <br />
              <span className="mt-2 block">
                You can override the default channel by specifying a channel name (e.g., #alerts) or a user (e.g., @username).
              </span>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slack_rich_message"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Message Format</FormLabel>
            <div className="flex items-center gap-2 mt-2">
              <FormControl>
                <Switch
                  checked={field.value || false}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <FormLabel className="text-sm font-normal">Send rich messages</FormLabel>
            </div>
            <FormDescription>
              Enable to send messages with rich formatting, attachments, and better visual presentation.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slack_channel_notify"
        render={({ field }) => (
          <FormItem>
            <div className="flex items-center gap-2">
              <FormControl>
                <Switch
                  checked={field.value || false}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <FormLabel>Notify Channel</FormLabel>
            </div>
            <FormDescription>
              When enabled, the message will trigger notifications for all channel members. Use with caution to avoid spam.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
    </>
  );
}
