#!/bin/sh
set -eu

API_URL=${API_URL:-}

cat >/usr/share/nginx/html/env.js <<EOF
/* generated each container start */
window.__CONFIG__ = {
  API_URL: "$API_URL"
};
EOF

# Continue with the official Nginx entrypoint
/docker-entrypoint.sh "$@"
