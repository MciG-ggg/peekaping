import type { ReactNode } from "react";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";
import { HomepageHeader } from "../../components/HomepageHeader";
import { KeyFeatures } from "../../components/KeyFeatures";

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Modern Uptime Monitoring`}
      description="A modern, self-hosted uptime monitoring solution. Monitor websites, APIs, and services with real-time notifications and beautiful status pages."
    >
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <KeyFeatures />
      </main>
    </Layout>
  );
}
