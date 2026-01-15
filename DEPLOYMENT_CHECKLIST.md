# Quick AWS EC2 Deployment Checklist

Complete these steps to deploy your app:

## ☐ Pre-Deployment (Local Machine)

- [ ] AWS CLI installed: `aws --version`
- [ ] AWS credentials configured: `aws sts get-caller-identity`
- [ ] Docker installed: `docker --version`
- [ ] Get AWS Account ID:
  ```bash
  aws sts get-caller-identity --query Account --output text
  ```

## ☐ AWS Setup

- [ ] EC2 instance running and Docker installed
- [ ] Know your EC2 Public IP address
- [ ] Have your EC2 key pair (.pem file)
- [ ] Create ECR repository:
  ```bash
  aws ecr create-repository --repository-name mongodb-user-app --region us-east-1
  ```

## ☐ Update Configuration

Edit `deploy.sh` with your details:
```bash
AWS_ACCOUNT_ID="YOUR_ID_HERE"
AWS_REGION="us-east-1"
EC2_INSTANCE_IP="YOUR_IP_HERE"
EC2_KEY_PATH="~/.ssh/your-key.pem"
```

## ☐ Build & Push to ECR (Local Machine)

```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project
chmod +x deploy.sh
./deploy.sh
```

Wait for: ✅ "Image pushed successfully!"

## ☐ Connect to EC2

```bash
ssh -i ~/.ssh/your-key.pem ec2-user@YOUR_EC2_IP
```

## ☐ Login to ECR (On EC2)

```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
```

## ☐ Pull Image (On EC2)

```bash
docker pull YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

## ☐ Run Container (On EC2)

```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://YOUR_MONGODB_HOST:27017" \
  YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

## ☐ Verify Running

```bash
docker ps
docker logs mongodb-app
curl http://localhost:8080/ping
```

## ☐ Update Security Group

AWS Console → EC2 → Security Groups:
- [ ] Add Inbound Rule: HTTP, Port 80, Source 0.0.0.0/0

## ☐ Test in Browser

Visit: `http://YOUR_EC2_PUBLIC_IP`

---

## Quick Command Reference

### If container won't start:
```bash
docker logs mongodb-app
docker rm mongodb-app  # Remove and try again
```

### To update the app:
```bash
# Local machine:
./deploy.sh

# On EC2:
docker stop mongodb-app
docker rm mongodb-app
docker pull YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
# Run the docker run command again
```

### View container logs live:
```bash
docker logs -f mongodb-app
```

---

**Estimated time:** 15-20 minutes

**Issues?** Check AWS_DEPLOYMENT.md for troubleshooting
