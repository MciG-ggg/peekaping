import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import Layout from "@/layout";

const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string().optional(),
});

type FormType = z.infer<typeof formSchema>;

const NewStatusPage = () => {
  const form = useForm<FormType>({
    defaultValues: {
      name: "",
      description: "",
    },
    resolver: zodResolver(formSchema),
  });

  const onSubmit = (data: FormType) => {
    // Submission logic will go here
    // For now, just log the data
    console.log(data);
  };

  return (
    <Layout pageName="New Status Page">
      <div className="flex flex-col gap-4">
        <p className="text-gray-500">
          Create a new status page to share your service status with users.
        </p>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-6 max-w-[600px]"
          >
            <div>
              <label className="block mb-1 font-medium" htmlFor="name">
                Name
              </label>
              <Input
                id="name"
                {...form.register("name")}
                placeholder="Status page name"
              />
              {form.formState.errors.name && (
                <p className="text-red-500 text-sm mt-1">
                  {form.formState.errors.name.message}
                </p>
              )}
            </div>
            <div>
              <label className="block mb-1 font-medium" htmlFor="description">
                Description
              </label>
              <Textarea
                id="description"
                {...form.register("description")}
                placeholder="Short description (optional)"
              />
              {form.formState.errors.description && (
                <p className="text-red-500 text-sm mt-1">
                  {form.formState.errors.description.message}
                </p>
              )}
            </div>
            <Button type="submit">Save</Button>
          </form>
        </Form>
      </div>
    </Layout>
  );
};

export default NewStatusPage;
