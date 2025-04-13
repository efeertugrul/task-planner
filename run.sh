#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting Todo Planning Application...${NC}"

# Initialize the database
echo -e "${GREEN}Initializing database...${NC}"
go run cmd/cli/main.go init-db
if [ $? -ne 0 ]; then
    echo "Failed to initialize database"
    exit 1
fi

# Start the Go server in the background
echo -e "${GREEN}Starting Go server...${NC}"
go run cmd/api/main.go &
SERVER_PID=$!

# Wait for the server to start
sleep 2

# Start the React web app
echo -e "${GREEN}Starting React web app...${NC}"
cd web
npm install
npm start &
WEB_PID=$!

# Function to handle cleanup
cleanup() {
    echo -e "${YELLOW}Shutting down...${NC}"
    kill $SERVER_PID
    kill $WEB_PID
    exit 0
}

# Trap SIGINT and SIGTERM signals and call cleanup
trap cleanup SIGINT SIGTERM

# Keep the script running
echo -e "${GREEN}Application is running!${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
wait 