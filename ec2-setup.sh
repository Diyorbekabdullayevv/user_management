#!/bin/bash

# Simple EC2 Deployment Script
# Run this on your EC2 instance after SSH

set -e

echo "🚀 MongoDB User App - EC2 Setup"
echo "================================"
echo ""

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get user input
read -p "Enter your AWS Account ID: " AWS_ACCOUNT_ID
read -p "Enter your AWS Region (default: us-east-1): " AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}

read -p "Enter MongoDB connection URI (mongodb://localhost:27017): " MONGO_URI
MONGO_URI=${MONGO_URI:-mongodb://localhost:27017}

echo ""
echo -e "${BLUE}1️⃣  Installing AWS CLI...${NC}"
# AWS CLI might already be installed
which aws > /dev/null 2>&1 || {
    echo "Installing AWS CLI..."
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
    rm -rf aws awscliv2.zip
}
echo -e "${GREEN}✅ AWS CLI ready${NC}"

echo ""
echo -e "${BLUE}2️⃣  Logging into AWS ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
echo -e "${GREEN}✅ Logged into ECR${NC}"

echo ""
echo -e "${BLUE}3️⃣  Pulling Docker image...${NC}"
docker pull $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/mongodb-user-app:latest
echo -e "${GREEN}✅ Image pulled${NC}"

echo ""
echo -e "${BLUE}4️⃣  Stopping old container (if exists)...${NC}"
docker stop mongodb-app 2>/dev/null || true
docker rm mongodb-app 2>/dev/null || true
echo -e "${GREEN}✅ Old container removed${NC}"

echo ""
echo -e "${BLUE}5️⃣  Starting new container...${NC}"
docker run -d \
  --name mongodb-app \
  -p 80:8080 \
  -e MONGO_URI="$MONGO_URI" \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/mongodb-user-app:latest

echo -e "${GREEN}✅ Container started${NC}"

echo ""
echo -e "${BLUE}6️⃣  Verifying...${NC}"
sleep 3

if docker ps | grep -q mongodb-app; then
    echo -e "${GREEN}✅ Container is running!${NC}"
    echo ""
    echo -e "${GREEN}🎉 SUCCESS! Your app is deployed${NC}"
    echo ""
    echo "📋 Quick commands:"
    echo "   View logs:     docker logs -f mongodb-app"
    echo "   Stop app:      docker stop mongodb-app"
    echo "   Restart app:   docker restart mongodb-app"
    echo "   Check status:  docker ps"
    echo ""
    echo "🌐 Access your app at: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"
else
    echo -e "${RED}❌ Container failed to start${NC}"
    echo ""
    echo "Debug info:"
    docker logs mongodb-app
    exit 1
fi
