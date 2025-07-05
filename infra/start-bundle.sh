#!/bin/bash

# Peekaping Bundle Docker Variants Launcher
# This script helps you easily start any of the three bundle variants

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_banner() {
    echo -e "${BLUE}"
    echo "╔═══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                         Peekaping Bundle Launcher                            ║"
    echo "║                                                                               ║"
    echo "║  Choose your preferred database backend for the bundle deployment            ║"
    echo "╚═══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_options() {
    echo -e "${YELLOW}Available variants:${NC}"
    echo ""
    echo -e "${GREEN}1. SQLite${NC} (Recommended for development)"
    echo "   • Simplest setup"
    echo "   • No external database required"
    echo "   • Best for: < 1000 concurrent users"
    echo "   • Memory usage: ~100-200MB"
    echo ""
    echo -e "${GREEN}2. MongoDB${NC} (NoSQL database)"
    echo "   • Full-featured document database"
    echo "   • Good for complex data structures"
    echo "   • Best for: 1000-10000 concurrent users"
    echo "   • Memory usage: ~200-500MB"
    echo ""
    echo -e "${GREEN}3. PostgreSQL${NC} (SQL database)"
    echo "   • Full-featured relational database"
    echo "   • Excellent for complex queries"
    echo "   • Best for: 1000-10000 concurrent users"
    echo "   • Memory usage: ~200-500MB"
    echo ""
}

run_variant() {
    local variant=$1
    local compose_file=""
    local variant_name=""

    case $variant in
        1)
            compose_file="docker-compose.bundle-sqlite.yml"
            variant_name="SQLite"
            ;;
        2)
            compose_file="docker-compose.bundle-mongo.yml"
            variant_name="MongoDB"
            ;;
        3)
            compose_file="docker-compose.bundle-postgres.yml"
            variant_name="PostgreSQL"
            ;;
        *)
            echo -e "${RED}Invalid option. Please choose 1, 2, or 3.${NC}"
            return 1
            ;;
    esac

    if [ ! -f "$compose_file" ]; then
        echo -e "${RED}Error: $compose_file not found!${NC}"
        echo "Please make sure you're running this script from the project root directory."
        return 1
    fi

    echo -e "${BLUE}Starting Peekaping with $variant_name...${NC}"
    echo ""

    # Check if Docker and Docker Compose are available
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Error: Docker is not installed or not in PATH${NC}"
        return 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}Error: Docker Compose is not installed or not in PATH${NC}"
        return 1
    fi

    # Stop any existing containers
    echo "Stopping any existing containers..."
    docker-compose -f "$compose_file" down 2>/dev/null || true

    # Start the selected variant
    echo "Starting $variant_name variant..."
    docker-compose -f "$compose_file" up -d

    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ Peekaping ($variant_name) started successfully!${NC}"
        echo ""
        echo -e "${YELLOW}Access your application at: http://localhost:8383${NC}"
        echo ""
        echo "Useful commands:"
        echo "  • View logs: docker-compose -f $compose_file logs -f"
        echo "  • Stop: docker-compose -f $compose_file down"
        echo "  • Restart: docker-compose -f $compose_file restart"
        echo ""
    else
        echo -e "${RED}✗ Failed to start Peekaping ($variant_name)${NC}"
        echo ""
        echo "Check the logs with: docker-compose -f $compose_file logs"
        return 1
    fi
}

main() {
    print_banner
    print_options

    echo -e "${YELLOW}Choose an option (1-3):${NC}"
    read -p "> " choice

    echo ""
    run_variant "$choice"
}

# Handle command line arguments
if [ $# -eq 1 ]; then
    case $1 in
        sqlite|1)
            run_variant 1
            ;;
        mongo|mongodb|2)
            run_variant 2
            ;;
        postgres|postgresql|3)
            run_variant 3
            ;;
        -h|--help)
            print_banner
            print_options
            echo -e "${YELLOW}Usage:${NC}"
            echo "  $0                    # Interactive mode"
            echo "  $0 sqlite             # Start SQLite variant"
            echo "  $0 mongo              # Start MongoDB variant"
            echo "  $0 postgres           # Start PostgreSQL variant"
            echo "  $0 --help             # Show this help"
            ;;
        *)
            echo -e "${RED}Invalid option: $1${NC}"
            echo "Use '$0 --help' for usage information."
            exit 1
            ;;
    esac
else
    main
fi
