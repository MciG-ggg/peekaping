import type { ReactNode } from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";
import Heading from "@theme/Heading";

import styles from "./index.module.css";

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx("hero ", styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.heroDescription}>
          <p>
            A modern, self-hosted uptime monitoring solution built with Go and
            React. Monitor your websites, APIs, and services with real-time
            notifications, beautiful status pages, and comprehensive analytics.
          </p>
        </div>
        <div className={styles.buttons}>
          <Link className="button button--primary button--lg" to="/intro">
            Get Started 🚀
          </Link>
          <Link
            className="button button--secondary button--lg"
                          to="/tutorial-basics/create-monitor"
            style={{ marginLeft: "1rem" }}
          >
            Quick Setup
          </Link>
        </div>
        <div className={styles.betaNotice}>
          <p>
            ⚠️ <strong>Beta Status:</strong> Peekaping is currently in beta.
            <Link to="/intro#️-beta-status"> Learn more</Link>
          </p>
        </div>
      </div>
    </header>
  );
}

function KeyFeatures() {
  return (
    <section className={styles.keyFeatures}>
      <div className="container">
        <div className="row">
          <div className="col col--12">
            <Heading as="h2">🚀 Quick Start with Docker</Heading>
            <p>
              Get Peekaping running in minutes with our Docker setup. No complex
              configuration required - just download and run!
            </p>

            <div className={styles.codeBlock}>
              <pre>
                <code>
                  {`# Download configuration files
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/.env.example -o .env
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/docker-compose.prod.yml -o docker-compose.yml

# Start Peekaping
docker compose up -d

# Visit http://localhost:8383`}
                </code>
              </pre>
            </div>
          </div>
        </div>

        <div className="row" style={{ marginTop: "2rem" }}>
          <div className="col col--12">
            <Heading as="h2">✨ Key Features</Heading>
            <ul className={styles.featureList}>
              <li>
                📊 <strong>Real-time Dashboard</strong> - Live status updates
                with WebSocket
              </li>
              <li>
                🔔 <strong>Smart Notifications</strong> - Email, Slack,
                Telegram, Webhooks
              </li>
              <li>
                📄 <strong>Public Status Pages</strong> - Share service status
                with users
              </li>
              <li>
                🛠 <strong>Maintenance Windows</strong> - Schedule maintenance to
                prevent false alerts
              </li>
              <li>
                🌐 <strong>Proxy Support</strong> - Route monitoring through
                HTTP proxies
              </li>
            </ul>
          </div>
        </div>
        <div className={styles.techStack}>
          <Heading as="h3">Built with Modern Technology</Heading>
          <div className={styles.badges}>
            <span className={styles.badge}>Go Backend</span>
            <span className={styles.badge}>React Frontend</span>
            <span className={styles.badge}>MongoDB</span>
            <span className={styles.badge}>Docker</span>
            <span className={styles.badge}>TypeScript</span>
            <span className={styles.badge}>WebSockets</span>
          </div>
        </div>
        <div className={styles.finalCta}>
          <Heading as="h3">Ready to start monitoring?</Heading>
          <p>Join the community and start monitoring your services today!</p>
          <div className={styles.buttons}>
            <Link className="button button--primary button--lg" to="/intro">
              Read Documentation
            </Link>
            <Link
              className="button button--outline button--lg"
              href="https://github.com/0xfurai/peekaping"
            >
              View on GitHub
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

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
