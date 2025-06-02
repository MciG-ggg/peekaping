import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import CreateEditProxy from "../components/create-edit-proxy";

const NewProxy = () => {
  const navigate = useNavigate();

  return (
    <Layout pageName="New Proxy">
      <CreateEditProxy onSuccess={() => navigate("/proxies")}/>
    </Layout>
  );
};

export default NewProxy;
