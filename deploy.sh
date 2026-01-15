#!/bin/bash

# MongoDB User Management - AWS Deployment Script
# This script builds and pushes your Docker image to AWS ECR

# ============== CONFIGURATION ==============
# Change these values with your AWS details

AWS_ACCOUNT_ID="739324486617"           # Get from AWS Console
AWS_REGION="eu-north-1"                          # Change if your region is different
ECR_REPO_NAME="mongodb-user-app"                # Name for your ECR repository
EC2_INSTANCE_IP="13.53.182.237"  # Get from AWS EC2 Dashboard
EC2_USER="ubuntu"                             # For Amazon Linux 2, use "ubuntu" for Ubuntu
EC2_KEY_PATH="~/.ssh/docker-key.pem"        # Path to your EC2 key pair

# ============== BUILD & PUSH ==============

echo "🔧 Building Docker image..."
docker build -t mongodb-user-app:latest .

echo "📝 Tagging image for AWS ECR..."
docker tag mongodb-user-app:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:latest

echo "🔐 Logging into AWS ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

echo "📤 Pushing image to ECR..."
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:latest

echo "✅ Image pushed successfully!"
echo ""
echo "Next steps:"
echo "1. SSH into your EC2 instance:"
echo "   ssh -i $EC2_KEY_PATH $EC2_USER@$EC2_INSTANCE_IP"
echo ""
echo "2. On your EC2 instance, run:"
echo "   aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"
echo "   docker run -d -p 80:8080 --name mongodb-app -e MONGO_URI='mongodb://YOUR_MONGODB_HOST:27017' $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME:latest"
