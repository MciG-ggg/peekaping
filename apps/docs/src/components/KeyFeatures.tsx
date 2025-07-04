import Link from "@docusaurus/Link";
import Heading from "@theme/Heading";

const code = `# Download configuration files
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/.env.example -o .env
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/docker-compose.prod.sqlite.yml -o docker-compose.yml

# Start Peekaping
docker compose up -d

# Visit http://localhost:8383`;

export function KeyFeatures() {
  return (
    <section className="py-16 lg:py-24 bg-white dark:bg-gray-900">
      <div className="container mx-auto px-6">
        {/* Quick Start Section */}
        <div className="max-w-4xl mx-auto mb-20">
          <div className="text-center mb-12">
            <Heading
              as="h2"
              className="text-3xl lg:text-4xl font-bold text-gray-900 dark:text-white mb-4"
            >
              🚀 Quick Start with Docker
            </Heading>
            <p className="text-lg text-gray-600 dark:text-gray-300 ">
              Get Peekaping running in minutes with our Docker setup. No complex
              configuration required - just download and run!
            </p>
          </div>

          <div className="bg-gray-900 dark:bg-gray-800 ">
            <pre className="text-sm lg:text-base text-gray-100 overflow-x-auto">
              <code>{code}</code>
            </pre>
          </div>
        </div>

        {/* Key Features Section */}
        <div className="max-w-4xl mx-auto mb-20">
          <div className="text-center mb-12">
            <Heading
              as="h2"
              className="text-3xl lg:text-4xl font-bold text-gray-900 dark:text-white mb-4"
            >
              ✨ Key Features
            </Heading>
          </div>
          <div className="grid gap-6 md:gap-8">
            <div className="grid md:grid-cols-2 gap-6">
              <div className="bg-blue-50 dark:bg-blue-900/20 p-6 rounded-xl border border-blue-200 dark:border-blue-800">
                <div className="flex items-center mb-3">
                  <span className="text-2xl mr-3">📊</span>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    Real-time Dashboard
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Live status updates with WebSocket
                </p>
              </div>

              <div className="bg-green-50 dark:bg-green-900/20 p-6 rounded-xl border border-green-200 dark:border-green-800">
                <div className="flex items-center mb-3">
                  <span className="text-2xl mr-3">🔔</span>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    Smart Notifications
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Email, Slack, Telegram, Webhooks
                </p>
              </div>

              <div className="bg-purple-50 dark:bg-purple-900/20 p-6 rounded-xl border border-purple-200 dark:border-purple-800">
                <div className="flex items-center mb-3">
                  <span className="text-2xl mr-3">📄</span>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    Public Status Pages
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Share service status with users
                </p>
              </div>

              <div className="bg-orange-50 dark:bg-orange-900/20 p-6 rounded-xl border border-orange-200 dark:border-orange-800">
                <div className="flex items-center mb-3">
                  <span className="text-2xl mr-3">🛠</span>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    Maintenance Windows
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Schedule maintenance to prevent false alerts
                </p>
              </div>
            </div>

            <div className="bg-indigo-50 dark:bg-indigo-900/20 p-6 rounded-xl border border-indigo-200 dark:border-indigo-800">
              <div className="flex items-center mb-3">
                <span className="text-2xl mr-3">🌐</span>
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Proxy Support
                </h3>
              </div>
              <p className="text-gray-600 dark:text-gray-300">
                Route monitoring through HTTP proxies
              </p>
            </div>
          </div>
        </div>

        {/* Technology Stack Section */}
        <div className="max-w-4xl mx-auto mb-20">
          <div className="text-center mb-12">
            <Heading
              as="h3"
              className="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-4"
            >
              Built with Modern Technology
            </Heading>
          </div>
          <div className="flex flex-wrap gap-3 justify-center">
            {[
              "Go Backend",
              "React Frontend",
              "MongoDB",
              "Postgres",
              "SQLite",
              "Docker",
              "TypeScript",
              "WebSockets",
            ].map((tech) => (
              <span
                key={tech}
                className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-2 rounded-full font-semibold text-sm shadow-lg hover:shadow-xl transition-shadow duration-200"
              >
                {tech}
              </span>
            ))}
          </div>
        </div>

        {/* Call to Action Section */}
        <div className="max-w-4xl mx-auto text-center">
          <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl p-8 lg:p-12 text-white">
            <Heading as="h3" className="text-2xl lg:text-3xl font-bold mb-4">
              Ready to start monitoring?
            </Heading>
            <p className="text-lg lg:text-xl mb-8 opacity-90">
              Join the community and start monitoring your services today!
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <Link
                className="button button--primary button--lg bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold text-lg hover:bg-gray-100 transition-colors duration-200 shadow-lg hover:shadow-xl"
                to="/intro"
              >
                Read Documentation
              </Link>
              <Link
                className="button button--outline button--lg bg-transparent text-white px-8 py-3 rounded-lg font-semibold text-lg border-2 border-white hover:bg-white hover:text-blue-600 transition-colors duration-200"
                href="https://github.com/0xfurai/peekaping"
              >
                View on GitHub
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
