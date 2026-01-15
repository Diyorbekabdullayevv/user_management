# 🚀 AWS EC2 Deployment - Complete Setup Guide

Everything is ready for deployment! Follow this guide to deploy your MongoDB User Management app to AWS EC2.

---

## 📊 Current Status

✅ Docker image built: `mongodb-user-app:latest` (64.4 MB)  
✅ Dockerfile updated for new project structure  
✅ Deployment scripts created  
✅ Documentation complete  

---

## 🎯 5-Minute Overview

```
Your Local Machine
    ↓ (docker build & push)
AWS ECR (Image Registry)
    ↓ (pull image)
EC2 Instance (Container Running)
    ↓ (port 80)
🌐 Your App: http://your-ec2-ip
```

---

## 📋 Pre-Deployment Checklist

- [ ] AWS Account created
- [ ] EC2 instance running with Docker installed
- [ ] AWS CLI installed: `aws --version`
- [ ] AWS credentials configured
- [ ] EC2 public IP address noted
- [ ] EC2 key pair (.pem) file available

---

## 🔧 Deployment Steps

### **Step 1: Get Your AWS Account ID** (2 min)

```bash
aws sts get-caller-identity --query Account --output text
```

**Output:** `123456789012`

Note this down!

---

### **Step 2: Get Your EC2 Public IP** (1 min)

Go to **AWS Console → EC2 → Instances**:
1. Select your instance
2. Copy the **Public IPv4** address
3. Note it down (example: `54.123.456.789`)

---

### **Step 3: Create ECR Repository** (1 min)

Run once in your local terminal:

```bash
aws ecr create-repository \
  --repository-name mongodb-user-app \
  --region us-east-1
```

Replace `us-east-1` if your region is different.

---

### **Step 4: Configure Deployment Script** (2 min)

Edit the `deploy.sh` file:

```bash
nano deploy.sh
```

Find and update these lines:

```bash
AWS_ACCOUNT_ID="123456789012"           # ← Update with your account ID
AWS_REGION="us-east-1"                   # ← Update if different region
ECR_REPO_NAME="mongodb-user-app"        # ← Keep as is
EC2_INSTANCE_IP="54.123.456.789"        # ← Update with your EC2 IP
EC2_USER="ec2-user"                     # ← Change to "ubuntu" for Ubuntu AMI
EC2_KEY_PATH="~/.ssh/your-key-pair.pem" # ← Update with your key path
```

Save with `CTRL+O`, `ENTER`, `CTRL+X`

---

### **Step 5: Build and Push to ECR** (3 min)

```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project
chmod +x deploy.sh
./deploy.sh
```

**Wait for:** ✅ "Image pushed successfully!"

This will:
- Build your Docker image
- Login to AWS ECR
- Push the image to your repository

---

### **Step 6: SSH into EC2** (1 min)

```bash
ssh -i ~/.ssh/your-key-pair.pem ec2-user@54.123.456.789
```

**For Ubuntu AMI:**
```bash
ssh -i ~/.ssh/your-key-pair.pem ubuntu@54.123.456.789
```

---

### **Step 7: Run Setup on EC2** (2 min)

**Option A: Easiest - Use Auto Setup Script**

```bash
curl -O https://raw.githubusercontent.com/YOUR_USERNAME/MongoDB-Project/main/ec2-setup.sh
chmod +x ec2-setup.sh
./ec2-setup.sh
```

Follow the prompts.

**Option B: Manual Setup**

```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | \
docker login --username AWS --password-stdin \
123456789012.dkr.ecr.us-east-1.amazonaws.com

# Pull the image
docker pull 123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest

# Run the container
docker run -d \
  -p 80:8080 \
  --name mongodb-app \
  -e MONGO_URI="mongodb://YOUR_MONGODB_HOST:27017" \
  123456789012.dkr.ecr.us-east-1.amazonaws.com/mongodb-user-app:latest
```

Replace:
- `123456789012` with your AWS Account ID
- `YOUR_MONGODB_HOST` with your MongoDB server IP/hostname

---

### **Step 8: Verify Container is Running** (1 min)

Still on EC2, run:

```bash
docker ps
```

Should show:
```
CONTAINER ID   IMAGE                    STATUS
abc123def456   .../mongodb-user-app    Up 2 seconds
```

---

### **Step 9: Update Security Group** (2 min)

Go to **AWS Console → EC2 → Security Groups**:

1. Find your instance's security group
2. Click **Edit Inbound Rules**
3. Click **Add Rule**
4. Set:
   - **Type:** HTTP
   - **Protocol:** TCP
   - **Port Range:** 80
   - **Source:** 0.0.0.0/0
5. **Save Rules**

---

### **Step 10: Access Your App** (1 min)

Open your browser and visit:

```
http://54.123.456.789
```

You should see your User Management app! 🎉

---

## 🧪 Testing Your Deployment

### On EC2, test the API:

```bash
# Test ping endpoint
curl http://localhost:8080/ping

# Should return:
# {"message":"pong"}
```

### Register a user (from browser):
1. Go to `http://your-ec2-ip`
2. Fill in the registration form
3. Submit
4. Check if you see success message

---

## 🐛 Troubleshooting

### **Container won't start?**
```bash
docker logs mongodb-app
```
Check for connection errors to MongoDB.

### **Can't access the app (connection timeout)?**
1. Verify security group allows port 80
2. Verify container is running: `docker ps`
3. Check logs: `docker logs mongodb-app`

### **MongoDB connection error?**
Verify MONGO_URI is correct:
```bash
docker exec mongodb-app cat /proc/1/environ | tr ',' '\n' | grep MONGO
```

### **Port 80 already in use?**
Use a different port:
```bash
docker run -d -p 8888:8080 ...
# Then access at http://your-ip:8888
```

---

## 📝 MongoDB Connection Strings

### Local MongoDB on EC2:
```
mongodb://host.docker.internal:27017
```

### MongoDB on another server:
```
mongodb://192.168.1.100:27017
```

### MongoDB Atlas (Cloud):
```
mongodb+srv://username:password@cluster.mongodb.net/database
```

---

## 🔄 Updating Your App

After making changes to your code:

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
# Run the docker run command again (from Step 7)
```

---

## 💰 Cost Optimization

To save AWS costs:

### Use EC2 Free Tier:
- t2.micro or t3.micro instances
- 750 hours/month free (eligible accounts)

### Monitor costs:
- AWS Billing Dashboard
- Set up cost alerts

### Clean up:
- Stop instance when not using: `aws ec2 stop-instances --instance-ids i-xxx`
- Or use AWS Systems Manager to auto-stop

---

## 📚 Additional Resources

- **Dockerfile:** `dockerfile` - Build configuration
- **Deploy Script:** `deploy.sh` - Automates ECR push
- **EC2 Script:** `ec2-setup.sh` - Automates EC2 setup
- **Detailed Guide:** `AWS_DEPLOYMENT.md` - Full documentation
- **Checklist:** `DEPLOYMENT_CHECKLIST.md` - Step verification

---

## 🎯 Next Steps

1. ✅ Complete Step 1-10 above
2. ✅ Verify app is running
3. ✅ Test user registration and filtering
4. ✅ Monitor logs: `docker logs -f mongodb-app`
5. ✅ Plan MongoDB backup strategy
6. ✅ Set up auto-restart (optional)

---

## 🔐 Security Best Practices

### Restrict Security Group:
Instead of `0.0.0.0/0`, use your IP:
```
Source: YOUR_IP/32
```

### Use Environment Variables:
Don't hardcode passwords:
```bash
docker run -d \
  -e MONGO_URI="$MONGO_URI" \
  ...
```

### Keep Docker Updated:
```bash
docker pull your-image:latest
docker system prune
```

---

## 📞 Support

If you encounter issues:

1. Check container logs: `docker logs mongodb-app`
2. Test MongoDB connection: `curl mongodb://host:27017`
3. Verify security groups
4. Check AWS CloudWatch logs
5. Review `AWS_DEPLOYMENT.md` for detailed troubleshooting

---

## ✅ Success Indicators

You've successfully deployed when:

✅ `docker ps` shows `mongodb-app` running  
✅ `curl http://localhost:8080/ping` returns `{"message":"pong"}`  
✅ Browser shows your app at `http://your-ec2-ip`  
✅ Form submission works and stores data  
✅ User filtering works  
✅ `docker logs -f mongodb-app` shows no errors  

---

**Total Setup Time:** ~30 minutes  
**Estimated AWS Monthly Cost:** $5-15 (Free Tier eligible)  
**Status:** ✅ Ready to Deploy!

Good luck! 🚀
