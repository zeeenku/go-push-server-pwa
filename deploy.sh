#!/bin/bash

# Step 1: Load environment variables from .env file
if [ -f ".env" ]; then
  export $(cat .env | xargs)
else
  echo ".env file not found!"
  exit 1
fi

# Step 2: Verify that PROJECT_DIR is set
if [ -z "$PROJECT_DIR" ]; then
  echo "PROJECT_DIR is not set in the .env file!"
  exit 1
fi

# Step 3: Navigate to the project directory
cd "$PROJECT_DIR" || exit

# Step 4: Install the Go modules required for the project
echo "Installing Go modules..."
go mod tidy
go mod download

# Step 5: Build the Go project (including main.go and send_notification.go)
echo "Building the Go project..."
go build -o push-server main.go
go build -o send_notification send_notification.go

# Step 9: Restart nginx (optional)
echo "Restarting nginx service..."
sudo systemctl restart nginx

echo "Deployment completed successfully!"
