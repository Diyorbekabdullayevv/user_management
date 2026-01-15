# AWS EC2 Deployment - Quick Start

Your Docker image is built and ready to deploy! Follow these steps:

## 📋 What You Have

✅ Docker image built: `mongodb-user-app:latest` (64.4 MB)
✅ Updated Dockerfile for new project structure
✅ Deployment scripts ready

---

## 🚀 Step-by-Step Deployment

### Step 1: Get Your AWS Account Details

```bash
# Get your AWS Account ID
aws sts get-caller-identity --query Account --output text

# Find your EC2 Public IP in AWS Console
# EC2 → Instances → Select your instance → Copy Public IPv4
```

### Step 2: Create ECR Repository (One-time only)

```bash
aws ecr create-repository \
  --repository-name mongodb-user-app \
  --region us-east-1
```

Replace `us-east-1` with your AWS region if different.

### Step 3: Edit Configuration

Update `deploy.sh` with your information:

```bash
# Edit the file
nano deploy.sh

# Change these:
AWS_ACCOUNT_ID="123456789012"           # Your AWS Account ID
AWS_REGION="us-east-1"                   # Your region
EC2_INSTANCE_IP="54.123.456.789"        # Your EC2 public IP
EC2_KEY_PATH="~/.ssh/your-key-pair.pem" # Your .pem file path
```

### Step 4: Build and Push Image to AWS ECR

```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project

# Make script executable
chmod +x deploy.sh

# Run the deployment script
./deploy.sh
```

**Expected output:**
```
🔧 Building Docker image...
...
✅ Image pushed successfully!

Next steps:
1. SSH into your EC2 instance:
   ssh -i ~/.ssh/your-key.pem ec2-user@YOUR_EC2_PUBLIC_IP
```

### Step 5: SSH into EC2 Instance

```bash
ssh -i ~/.ssh/your-key-pair.pem ec2-user@YOUR_EC2_PUBLIC_IP
```

For Ubuntu AMI, use `ubuntu` instead of `ec2-user`:
```bash
ssh -i ~/.ssh/your-key-pair.pem ubuntu@YOUR_EC2_PUBLIC_IP
```

### Step 6: Run Setup Script on EC2 (Easiest Way)

Copy and run this on your EC2 instance:

```bash
curl -O https://raw.githubusercontent.com/YOUR_REPO/MongoDB-Project/main/ec2-setup.sh
chmod +x ec2-setup.sh
./ec2-setup.sh
```

**Or manually** (see "Manual Steps" below)

---

## Manual Steps (If Not Using Auto Script)

### On EC2 - Login to AWS ECR:
```bash
aws ecr get-login-password --region us-east-1 | \
docker login --username AWS --password-stdin \
123456789012.dkr.ecr.us-east-1.amazonaws.com
```

### On EC2 - Pull Image:
```bash
docker pull 123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

### On EC2 - Run Container:

**If MongoDB is on EC2 (local):**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://host.docker.internal:27017" \
  123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

**If MongoDB is on another server:**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://mongodb-server-ip:27017" \
  123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

**If using MongoDB Atlas (Cloud):**
```bash
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb+srv://user:pass@cluster.mongodb.net/dbname" \
  123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

---

## ✅ Verify It's Running

### On EC2, check container:
```bash
docker ps
```

Should show `mongodb-app` as running.

### View logs:
```bash
docker logs -f mongodb-app
```

### Test the API:
```bash
curl http://localhost:8080/ping
```

Should return: `{"message":"pong"}`

---

## 🌐 Access Your App from Browser

### Update Security Group First:

1. Go to **AWS Console → EC2 → Security Groups**
2. Find your instance's security group
3. Click **Edit Inbound Rules**
4. Add new rule:
   - Type: **HTTP**
   - Protocol: **TCP**
   - Port Range: **80**
   - Source: **0.0.0.0/0** (or your IP for security)
5. **Save**

### Visit Your App:
```
http://YOUR_EC2_PUBLIC_IP
```

Example: `http://54.123.456.789`

---

## 📊 Common Commands

| Task | Command |
|------|---------|
| View running containers | `docker ps` |
| View all containers | `docker ps -a` |
| View logs | `docker logs mongodb-app` |
| Live logs | `docker logs -f mongodb-app` |
| Stop container | `docker stop mongodb-app` |
| Start container | `docker start mongodb-app` |
| Restart container | `docker restart mongodb-app` |
| Remove container | `docker rm mongodb-app` |
| SSH into container | `docker exec -it mongodb-app sh` |

---

## 🔄 Update Your App (After Code Changes)

### On local machine:
```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project
./deploy.sh
```

### On EC2:
```bash
docker stop mongodb-app
docker rm mongodb-app
docker pull 123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
docker run -d -p 80:8080 --name mongodb-app \
  -e MONGO_URI="mongodb://..." \
  123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

---

## 🐛 Troubleshooting

### Container won't start?
```bash
docker logs mongodb-app
```

### Port 80 already in use?
```bash
# Use a different port
docker run -d -p 8080:8080 ...
# Then access at http://your-ip:8080
```

### Can't reach MongoDB?
- Check MONGO_URI is correct
- Verify MongoDB is running and accessible
- Check security groups allow port 27017

### Permission denied errors?
```bash
# Make scripts executable
chmod +x deploy.sh ec2-setup.sh
```

---

## 📁 Files Reference

| File | Purpose |
|------|---------|
| `dockerfile` | Docker configuration |
| `deploy.sh` | Build & push to ECR |
| `ec2-setup.sh` | Auto-setup script for EC2 |
| `AWS_DEPLOYMENT.md` | Detailed deployment guide |
| `DEPLOYMENT_CHECKLIST.md` | Step-by-step checklist |

---

## 🎯 Summary

1. ✅ Docker image built successfully
2. ⏳ Update `deploy.sh` with your AWS details
3. ⏳ Run `./deploy.sh` to push to ECR
4. ⏳ SSH into EC2 and run `ec2-setup.sh`
5. ⏳ Update Security Group to allow HTTP (port 80)
6. ✅ Visit `http://YOUR_EC2_PUBLIC_IP`

**Total time:** ~15-20 minutes

Need help? Check `AWS_DEPLOYMENT.md` for detailed instructions.
