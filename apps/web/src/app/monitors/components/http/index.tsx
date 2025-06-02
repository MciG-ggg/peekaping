import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useFormContext } from "react-hook-form";
import { z } from "zod";
import Advanced, { advancedDefaultValues, advancedSchema } from "./advanced";
import Authentication, {
  authenticationDefaultValues,
  authenticationSchema,
} from "./authentication";
import HttpOptions, {
  httpOptionsDefaultValues,
  httpOptionsSchema,
} from "./options";
import { Separator } from "@/components/ui/separator";

export const httpSchema = z
  .object({
    url: z.string().url({ message: "Invalid URL" }),
  })
  .and(advancedSchema)
  .and(httpOptionsSchema)
  .and(authenticationSchema);

type HttpForm = z.infer<typeof httpSchema>;

export const httpDefaultValues: HttpForm = {
  url: "https://example.com",

  ...httpOptionsDefaultValues,
  ...advancedDefaultValues,
  ...authenticationDefaultValues,
};

export const Essentials = () => {
  const form = useFormContext();

  return (
    <FormField
      control={form.control}
      name="http.url"
      render={({ field }) => (
        <FormItem>
          <FormLabel>URL</FormLabel>
          <FormControl>
            <Input placeholder="https://" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};

export const Other = () => {
  return (
    <>
      <Advanced />
      <Separator className="my-8" />
      <Authentication />
      <Separator className="my-8" />
      <HttpOptions />
    </>
  );
};
