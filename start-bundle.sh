#!/bin/bash

# Wrapper script for the Bundle Docker Variants Launcher
# This script calls the actual launcher located in the infra folder

set -e

# Change to the infra directory and run the actual script
cd "$(dirname "$0")/infra"
exec ./start-bundle.sh "$@"
