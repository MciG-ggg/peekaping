import type { ReactNode } from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Heading from "@theme/Heading";

export function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx("hero bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800")}>
      <div className="container mx-auto px-6 py-16 lg:py-24">
        <div className="text-center  mx-auto">
          <Heading as="h1" className="text-4xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-6">
            {siteConfig.title}
          </Heading>
          {/* <p className="text-xl lg:text-2xl text-gray-600 dark:text-gray-300 mb-8 font-medium">
            {siteConfig.tagline}
          </p> */}
          <div className="mb-10">
            <p className="text-lg lg:text-xl text-gray-700 dark:text-gray-300 leading-relaxed  mx-auto">
              A modern, self-hosted uptime monitoring solution built with Go and
              React. Monitor your websites, APIs, and services with real-time
              notifications, beautiful status pages, and comprehensive analytics.
            </p>
          </div>
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-8">
            <Link
              className="button button--primary button--lg bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-lg font-semibold text-lg transition-colors duration-200 shadow-lg hover:shadow-xl"
              to="/intro"
            >
              Docs 📚
            </Link>
            {/* <Link
              className="button button--secondary button--lg bg-white dark:bg-gray-800 text-blue-600 dark:text-blue-400 px-8 py-3 rounded-lg font-semibold text-lg border-2 border-blue-600 dark:border-blue-400 hover:bg-blue-50 dark:hover:bg-gray-700 transition-colors duration-200"
              to="/tutorial-basics/create-monitor"
            >
              Quick Setup
            </Link> */}
            <Link
              className="button button--outline button--lg bg-gradient-to-r from-green-500 to-blue-500 hover:from-green-600 hover:to-blue-600 text-white px-8 py-3 rounded-lg font-semibold text-lg transition-all duration-200 shadow-lg hover:shadow-xl"
              href="https://demo.peekaping.com"
              target="_blank"
              rel="noopener noreferrer"
            >
              Live Demo ✨
            </Link>
          </div>
          <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-4 inline-block">
            <div className="text-amber-800 dark:text-amber-200 font-medium">
              ⚠️ <strong>Beta Status:</strong> Peekaping is currently in beta.{" "}
              <Link
                to="/intro#️-beta-status"
                className="text-amber-900 dark:text-amber-100 underline hover:no-underline"
              >
                Learn more
              </Link>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
