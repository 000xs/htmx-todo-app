#!/bin/bash
set -e

# Configuration
APP_NAME="my-go-app"
GO_VERSION="1.23.0"
NGINX_PORT=80
APP_PORT=8080

# Update system
sudo apt-get update -y && sudo apt-get upgrade -y

# Install dependencies
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    ufw

# Install Docker
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
https://download.docker.com/linux/debian \
$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install Docker Compose
DOCKER_COMPOSE_VERSION="v2.27.0"
sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create Dockerfile
cat > Dockerfile <<EOF
# Build stage
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /${APP_NAME}

# Final stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /${APP_NAME} .
EXPOSE ${APP_PORT}
CMD ["./${APP_NAME}"]
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
      - "${NGINX_PORT}:80"
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
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    sendfile on;
    keepalive_timeout 65;

    upstream go-app {
        server go-app:${APP_PORT};
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
            
            proxy_connect_timeout 300;
            proxy_send_timeout 300;
            proxy_read_timeout 300;
            send_timeout 300;
        }
    }
}
EOF

# Configure firewall
sudo ufw allow ${NGINX_PORT}/tcp
sudo ufw --force enable

# Build and start containers
sudo docker-compose up -d --build

echo "Deployment completed!"
echo "Access your application at: http://$(curl -4 -s ifconfig.me)"