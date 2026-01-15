# AWS EC2 Deployment Guide

Deploy your MongoDB User Management app to AWS EC2 using Docker.

## Prerequisites

✅ AWS Account with EC2 instance running
✅ Docker installed on EC2 instance
✅ AWS CLI configured locally
✅ EC2 key pair (.pem file)

---

## Step 1: Get Your AWS Details

### Find your AWS Account ID:
```bash
aws sts get-caller-identity --query Account --output text
```

### Find your EC2 Instance:
1. Go to AWS Console → EC2 Dashboard
2. Note your:
   - **Public IPv4 address** (or DNS name)
   - **Key pair name** (.pem file location)

---

## Step 2: Update the Deployment Script

Edit `deploy.sh` and replace:

```bash
AWS_ACCOUNT_ID="123456789012"           # Your AWS Account ID
AWS_REGION="eu-north-1a"                   # Your AWS region
EC2_INSTANCE_IP="54.123.456.789"        # Your EC2 public IP
EC2_KEY_PATH="~/.ssh/your-key-pair.pem" # Path to your .pem file
```

---

## Step 3: Create ECR Repository (One-time)

```bash
aws ecr create-repository \
  --repository-name mongodb-user-app \
  --region us-east-1
```

---

## Step 4: Build and Push to ECR

```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project
chmod +x deploy.sh
./deploy.sh
```

This will:
1. Build the Docker image locally
2. Tag it for AWS ECR
3. Login to AWS ECR
4. Push the image to your repository

---

## Step 5: SSH into EC2 Instance

```bash
ssh -i ~/.ssh/your-key-pair.pem ec2-user@YOUR_EC2_PUBLIC_IP
# or for Ubuntu AMI:
# ssh -i ~/.ssh/your-key-pair.pem ubuntu@YOUR_EC2_PUBLIC_IP
```

---

## Step 6: Pull and Run the Image on EC2

### Login to ECR (on EC2):
```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
```

### Pull the image:
```bash
docker pull YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

### Run the container:

**Option A: If MongoDB is on EC2:**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://localhost:27017" \
  YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

**Option B: If MongoDB is on Atlas (Cloud):**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb+srv://username:password@cluster.mongodb.net/dbname" \
  YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

**Option C: If MongoDB is on another server:**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://mongodb-server-ip:27017" \
  YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

---

## Step 7: Verify the App is Running

### Check container status:
```bash
docker ps
```

### View logs:
```bash
docker logs mongodb-app
```

### Test the app:
```bash
curl http://localhost:8080/ping
```

---

## Step 8: Update EC2 Security Group

Your app is running on port 8080, but we mapped it to port 80.

### Open AWS Console:
1. Go to EC2 → Security Groups
2. Find your instance's security group
3. Edit Inbound Rules
4. Add rule: **Type: HTTP, Port: 80, Source: 0.0.0.0/0**

### Now access your app:
Open browser: `http://YOUR_EC2_PUBLIC_IP`

---

## Additional Commands

### Stop the container:
```bash
docker stop mongodb-app
```

### Remove the container:
```bash
docker rm mongodb-app
```

### View logs in real-time:
```bash
docker logs -f mongodb-app
```

### Restart the container:
```bash
docker restart mongodb-app
```

### Update the image:
```bash
# On your local machine:
./deploy.sh

# On EC2:
docker pull YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
docker stop mongodb-app
docker rm mongodb-app
# Run again with the docker run command above
```

---

## Environment Variables (Optional)

Update the Dockerfile to support environment variables for MongoDB connection:

```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_HOST="mongodb.example.com" \
  -e MONGO_PORT="27017" \
  -e MONGO_DB="users" \
  YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

---

## Troubleshooting

### Container won't start?
```bash
docker logs mongodb-app
```

### Port already in use?
```bash
# Change the port mapping:
docker run -d -p 8080:8080 ...  # Use 8080 instead of 80
```

### Can't connect to MongoDB?
- Ensure MongoDB is running and accessible
- Check security groups allow MongoDB port (27017)
- Verify connection string is correct

### Need to rebuild?
```bash
docker stop mongodb-app
docker rm mongodb-app
./deploy.sh
# Then run docker pull and docker run again
```

---

## Architecture

```
Local Machine
    ↓ (build & push)
AWS ECR (Image Registry)
    ↓ (pull & run)
EC2 Instance (Running Container)
    ↓ (port 80)
Browser
```

---

## Quick Reference

| Task | Command |
|------|---------|
| Build image | `docker build -t mongodb-user-app:latest .` |
| List images | `docker images` |
| Run container | `docker run -d -p 80:8080 ...` |
| Stop container | `docker stop mongodb-app` |
| View logs | `docker logs -f mongodb-app` |
| Exec into container | `docker exec -it mongodb-app sh` |

---

## Support

For more help:
- AWS ECR: https://docs.aws.amazon.com/ecr/
- Docker: https://docs.docker.com/
- EC2: https://docs.aws.amazon.com/ec2/
