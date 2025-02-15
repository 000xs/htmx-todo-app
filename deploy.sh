#!/bin/bash
set -e

# Configuration
APP_NAME="my-go-app"
GO_VERSION="1.23.0"
NGINX_PORT=80
APP_PORT=8080

# Update system
sudo apt-get update -y && sudo apt-get upgrade -y

# Install Docker
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.27.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create Dockerfile for Go app
cat > Dockerfile <<EOF
# Build stage
FROM golang:$GO_VERSION-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /$APP_NAME

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /$APP_NAME .
EXPOSE $APP_PORT
CMD ["./$APP_NAME"]
EOF

# Create docker-compose.yml
cat > docker-compose.yml <<EOF
version: '3.8'

services:
  go-app:
    build: .
    restart: always
    environment:
      - APP_ENV=production
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    restart: always
    ports:
      - "$NGINX_PORT:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - go-app
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
EOF

# Create Nginx configuration
cat > nginx.conf <<EOF
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    upstream go-app {
        server go-app:$APP_PORT;
    }

    server {
        listen 80;
        server_name _;

        location / {
            proxy_pass http://go-app;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }
    }
}
EOF

# Build and start containers
sudo docker-compose up -d --build

# Configure firewall
sudo ufw allow $NGINX_PORT/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable

echo "Deployment completed successfully!"
echo "Application available at: http://$(curl -s ifconfig.me)"