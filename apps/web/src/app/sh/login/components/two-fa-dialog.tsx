import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import {
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  Form,
} from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

interface TwoFADialogProps {
  email: string;
  password: string;
  onSubmit: (data: { email: string; password: string; code: string }) => void;
  error?: string | null;
  loading?: boolean;
}

const twoFASchema = z.object({
  code: z.string().min(1, "2FA code is required"),
});

type TwoFAFormValues = z.infer<typeof twoFASchema>;

export function TwoFADialog({
  email,
  password,
  onSubmit,
  error,
  loading,
}: TwoFADialogProps) {
  const form = useForm<TwoFAFormValues>({
    resolver: zodResolver(twoFASchema),
    defaultValues: { code: "" },
  });

  function handleSubmit(values: TwoFAFormValues) {
    onSubmit({ email, password, code: values.code });
  }

  return (
    <Card>
      <CardHeader className="text-center">
        <CardTitle className="text-xl">
          Two-Factor Authentication Required
        </CardTitle>
        <CardDescription>
          2FA is enabled for your account. Please enter your authentication code
          to continue.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="grid gap-6"
          >
            <FormItem>
              <FormLabel>2FA Code</FormLabel>
              <FormControl>
                <Input
                  placeholder="Enter 2FA code"
                  {...form.register("code")}
                  autoFocus
                />
              </FormControl>
              <FormMessage />
            </FormItem>

            {error && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Verifying..." : "Verify 2FA"}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
