import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { TypographyH4 } from "@/components/ui/typography";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useFormContext, useWatch } from "react-hook-form";
import { z } from "zod";

const responseValidationTypes = [
  { value: "none", label: "No validation" },
  { value: "keyword", label: "Keyword Search" },
  { value: "json_query", label: "JSON Query Expression" },
];

const jsonQueryConditions = [
  { value: "===", label: "Equal (===)" },
  { value: "!=", label: "Not Equal (!=)" },
  { value: ">", label: "Greater Than (>)" },
  { value: "<", label: "Less Than (<)" },
  { value: ">=", label: "Greater Than or Equal (>=)" },
  { value: "<=", label: "Less Than or Equal (<=)" },
];

const ResponseValidation = () => {
  const form = useFormContext();
  const responseValidation = useWatch({
    control: form.control,
    name: "response_validation",
  });

  return (
    <>
      <TypographyH4>Response Validation</TypographyH4>
      
      <FormField
        control={form.control}
        name="response_validation"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Validation Type</FormLabel>
            <Select
              onValueChange={(value) => {
                field.onChange(value);
                // Reset related fields when validation type changes
                if (value === "none") {
                  form.setValue("keyword", "");
                  form.setValue("invert_keyword", false);
                  form.setValue("json_query", "");
                  form.setValue("json_query_condition", "===");
                  form.setValue("json_query_expected_value", "");
                }
              }}
              value={field.value}
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select validation type" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {responseValidationTypes.map((type) => (
                  <SelectItem key={type.value} value={type.value}>
                    {type.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>
              Choose how to validate the HTTP response content.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {responseValidation === "keyword" && (
        <>
          <FormField
            control={form.control}
            name="keyword"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Keyword</FormLabel>
                <FormControl>
                  <Input placeholder="Enter keyword to search for" {...field} />
                </FormControl>
                <FormDescription>
                  Search keyword in plain HTML or JSON response. The search is case-sensitive.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="invert_keyword"
            render={({ field }) => (
              <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                <FormControl>
                  <Checkbox
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
                <div className="space-y-1 leading-none">
                  <FormLabel>Invert Keyword</FormLabel>
                  <FormDescription>
                    Look for the keyword to be absent rather than present.
                  </FormDescription>
                </div>
              </FormItem>
            )}
          />
        </>
      )}

      {responseValidation === "json_query" && (
        <>
          <FormField
            control={form.control}
            name="json_query"
            render={({ field }) => (
              <FormItem>
                <FormLabel>JSON Query Expression</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="$"
                    {...field}
                    rows={3}
                  />
                </FormControl>
                <FormDescription>
                  Parse and extract specific data from the server's JSON response using a simple JSON path
                  or use "$" for the raw response, if not expecting JSON. The result is then
                  compared to the expected value, as strings. Supports basic dot notation like "user.name" 
                  or "data.items[0].id" for simple JSON path queries.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              control={form.control}
              name="json_query_condition"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Condition</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    value={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select condition" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {jsonQueryConditions.map((condition) => (
                        <SelectItem key={condition.value} value={condition.value}>
                          {condition.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="json_query_expected_value"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Expected Value</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter expected value" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
        </>
      )}
    </>
  );
};

export default ResponseValidation;