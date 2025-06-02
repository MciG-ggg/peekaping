import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import CreateEditNotifier from "../create-edit-notifier";

const NewNotifier = () => {
  const navigate = useNavigate();

  return (
    <Layout pageName="New Notifier">
      <CreateEditNotifier onSuccess={() => navigate("/notifiers")} />
    </Layout>
  );
};

export default NewNotifier;
