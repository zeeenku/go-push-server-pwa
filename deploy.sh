#!/bin/bash

# Load environment variables from .env file (with handling for spaces and special characters)
if [ -f "$PROJECT_DIR/.env" ]; then
  export $(grep -v '^#' $PROJECT_DIR/.env | xargs -d '\n')
else
  echo ".env file not found!"
  exit 1
fi

# Step 1: Pull latest changes from the repository
echo "Pulling latest changes from Git..."
cd $PROJECT_DIR
git pull

# Step 2: Install dependencies (if needed)
# For Go, you might want to download new dependencies
# go get -u ./...

# Step 3: Build the Go app (if using a compiled binary)
echo "Building the Go project..."
go build -o push-server main.go

# Step 4: Create or update the systemd service file
echo "Creating systemd service file..."
cat > $SERVICE_FILE <<EOF
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

# Step 5: Reload systemd to apply the new service file
echo "Reloading systemd to apply new service..."
sudo systemctl daemon-reload

# Step 6: Enable and start the service
echo "Enabling and starting the push-server service..."
sudo systemctl enable push-server.service
sudo systemctl start push-server.service

# Step 7: Restart nginx (optional)
echo "Restarting nginx service..."
sudo systemctl restart nginx

echo "Deployment completed successfully!"
