#!/bin/bash

# TicPin Database Seeding Script
# This script helps manage seeding test data into Firebase

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}    TicPin Database Management${NC}"
echo -e "${GREEN}=====================================${NC}\n"

# Check if we're in the right directory
if [ ! -f "$BACKEND_DIR/go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this from Backend(Go) directory${NC}"
    exit 1
fi

# Show menu
if [ $# -eq 0 ]; then
    echo "Usage: ./seed.sh [command]"
    echo ""
    echo "Commands:"
    echo "  seed      - Seed database with test data (plays, dining, events)"
    echo "  clear     - Clear all data from database (⚠️ DESTRUCTIVE)"
    echo "  reset     - Clear all data and re-seed with test data"
    echo ""
    echo "Examples:"
    echo "  ./seed.sh seed      # Add test data"
    echo "  ./seed.sh clear     # Remove all data"
    echo "  ./seed.sh reset     # Start fresh with new test data"
    exit 0
fi

COMMAND=$1

case $COMMAND in
    seed)
        echo -e "${YELLOW}Seeding database with test data...${NC}\n"
        cd "$BACKEND_DIR"
        go run scripts/seed/main.go
        echo -e "\n${GREEN}✓ Database seeding complete!${NC}"
        echo -e "Check Firebase Console to verify data"
        ;;
    
    clear)
        echo -e "${RED}⚠️  WARNING: This will DELETE ALL data from the database!${NC}"
        echo -e "${RED}⚠️  This action CANNOT be undone.${NC}\n"
        read -p "Are you sure? Type 'yes' to confirm: " confirmation
        
        if [ "$confirmation" != "yes" ]; then
            echo -e "${YELLOW}Cancelled.${NC}"
            exit 0
        fi
        
        echo -e "${YELLOW}Clearing database...${NC}\n"
        cd "$BACKEND_DIR"
        go run scripts/clear/main.go
        echo -e "\n${GREEN}✓ Database cleared!${NC}"
        ;;
    
    reset)
        echo -e "${RED}⚠️  WARNING: This will DELETE ALL data and reseed!${NC}"
        read -p "Are you sure? Type 'yes' to confirm: " confirmation
        
        if [ "$confirmation" != "yes" ]; then
            echo -e "${YELLOW}Cancelled.${NC}"
            exit 0
        fi
        
        echo -e "${YELLOW}Clearing database...${NC}"
        cd "$BACKEND_DIR"
        go run scripts/clear/main.go
        
        echo -e "\n${YELLOW}Reseeding database with test data...${NC}\n"
        go run scripts/seed/main.go
        
        echo -e "\n${GREEN}✓ Database reset and reseeded!${NC}"
        ;;
    
    *)
        echo -e "${RED}Unknown command: $COMMAND${NC}"
        echo "Run './seed.sh' with no arguments for help"
        exit 1
        ;;
esac
