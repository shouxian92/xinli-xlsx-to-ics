# Render Deployment Guide

## Prerequisites

1. **GitHub Account**: Your code must be in a GitHub repository
2. **Render Account**: Sign up at [render.com](https://render.com)
3. **Telegram Bot Token**: Get from [@BotFather](https://t.me/botfather)
4. **Linux/macOS Environment**: For local development (Windows users can use WSL)

## Deployment Steps

### 1. Push Code to GitHub

```bash
# Add all files (excluding Windows-specific ones)
git add .
git commit -m "Add Render deployment files and Linux support"
git push origin main
```

### 2. Connect to Render

1. Go to [render.com](https://render.com) and sign in
2. Click "New +" → "Web Service"
3. Connect your GitHub account if not already connected
4. Select your `xinli-xlsx-to-ics` repository

### 3. Configure the Service

- **Name**: `xlsx-to-ics-bot` (or any name you prefer)
- **Environment**: `Docker`
- **Region**: Choose closest to your users
- **Branch**: `main`
- **Build Command**: Leave empty (uses Dockerfile)
- **Start Command**: Leave empty (uses Dockerfile)

### 4. Set Environment Variables

Click "Environment" tab and add:

| Key | Value | Description |
|-----|-------|-------------|
| `TELEGRAM_BOT_TOKEN` | `your_bot_token_here` | Your actual bot token |
| `PORT` | `8080` | Port for health checks |

### 5. Deploy

1. Click "Create Web Service"
2. Render will automatically build and deploy your bot
3. Wait for build to complete (usually 2-5 minutes)

### 6. Test Your Bot

1. Copy the service URL from Render dashboard
2. Test the health endpoint: `https://your-service.onrender.com/health`
3. Start a conversation with your bot on Telegram
4. Send `/start` command
5. Upload an Excel file to test

## Important Notes

- **Free Tier**: 750 hours/month (enough for 24/7 operation)
- **Auto-sleep**: Free tier services sleep after 15 minutes of inactivity
- **Cold Start**: First request after sleep may take 10-30 seconds
- **Environment Variables**: Never commit real tokens to Git

## Troubleshooting

### Build Failures
- Check Dockerfile syntax
- Ensure all dependencies are in go.mod
- Verify Go version compatibility

### Bot Not Responding
- Check Render logs for errors
- Verify TELEGRAM_BOT_TOKEN is set correctly
- Ensure bot is not blocked by users

### Health Check Failures
- Verify `/health` endpoint is working
- Check if PORT environment variable is set
- Review Render service logs

## Cost Optimization

Since your bot runs every 6 months:
- **Option 1**: Keep running (stays within free tier)
- **Option 2**: Destroy service when not needed, redeploy when needed
- **Option 3**: Use Render's auto-sleep feature (free tier)

## Monitoring

- **Logs**: Available in Render dashboard
- **Health**: Automatic health checks every 30 seconds
- **Uptime**: Render provides uptime monitoring
- **Scaling**: Can easily upgrade to paid plans if needed
