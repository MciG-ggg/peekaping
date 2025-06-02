import Layout from "@/layout";
import { useNavigate } from "react-router-dom";

const StatusPagesPage = () => {
  const navigate = useNavigate();

  return (
    <Layout
      pageName="Status pages"
      onCreate={() => {
        navigate("/status-pages/new");
      }}
    >
      <div>StatusPagesPage</div>
    </Layout>
  );
};

export default StatusPagesPage;
