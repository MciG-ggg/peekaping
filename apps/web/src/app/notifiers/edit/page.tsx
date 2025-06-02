import { useParams, useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getNotificationsByIdOptions,
  getNotificationsByIdQueryKey,
  putNotificationsByIdMutation,
} from "@/api/@tanstack/react-query.gen";
import Layout from "@/layout";
import CreateEditNotifier from "../create-edit-notifier";
import { toast } from "sonner";

const EditNotifier = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    ...getNotificationsByIdOptions({ path: { id: id! } }),
    enabled: !!id,
  });

  const mutation = useMutation({
    ...putNotificationsByIdMutation(),
    onSuccess: () => {
      toast.success("Notifier updated successfully");
      queryClient.removeQueries({
        queryKey: getNotificationsByIdQueryKey({ path: { id: id! } }),
      });
      navigate("/notifiers");
    },
    onError: () => {
      toast.error("Failed to update notifier");
    },
  });

  if (isLoading) return <Layout pageName="Edit Notifier">Loading...</Layout>;
  if (error || !data?.data)
    return <Layout pageName="Edit Notifier">Error loading notifier</Layout>;

  // Prepare initial values for the form
  const notifier = data.data;
  const config = JSON.parse(notifier.config || "{}");

  const initialValues = {
    name: notifier.name || "",
    type: notifier.type,
    ...(config || {}),
  };

  const handleSubmit = (values: {
    name: string;
    type: string;
    config: string;
  }) => {
    mutation.mutate({
      path: { id: id! },
      body: {
        name: values.name,
        type: values.type,
        config: values.config,
        active: notifier.active,
        is_default: notifier.is_default,
      },
    });
  };

  return (
    <Layout pageName={`Edit Notifier: ${notifier.name}`}>
      <CreateEditNotifier
        initialValues={initialValues}
        onSubmit={handleSubmit}
        isLoading={mutation.isPending}
        isEdit
      />
    </Layout>
  );
};

export default EditNotifier;
