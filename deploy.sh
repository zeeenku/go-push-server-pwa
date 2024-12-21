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

# Step 5: Build the Go project
echo "Building the Go project..."
go build -o push-server main.go

# Step 6: Create or update the systemd service file
echo "Creating systemd service file..."
SERVICE_FILE="/etc/systemd/system/push-server.service"

cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Push Notification Server
After=network.target

[Service]
User=your-username
Group=your-username
WorkingDirectory=$PROJECT_DIR
ExecStart=$PROJECT_DIR/push-server
Restart=always
EnvironmentFile=$PROJECT_DIR/.env
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# Step 7: Reload systemd to apply the new service file
echo "Reloading systemd to apply new service..."
sudo systemctl daemon-reload

# Step 8: Enable and start the push-server service
echo "Enabling and starting the push-server service..."
sudo systemctl enable push-server.service
sudo systemctl start push-server.service

# Step 9: Restart nginx (optional)
echo "Restarting nginx service..."
sudo systemctl restart nginx

echo "Deployment completed successfully!"
