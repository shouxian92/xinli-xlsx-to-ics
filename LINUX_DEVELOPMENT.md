# Linux Development Guide

This guide is specifically for Linux and macOS users developing the XLSX to ICS Telegram Bot.

## 🐧 **Linux Setup**

### Prerequisites

- **Go 1.22.0+**: Install from [golang.org](https://golang.org/dl/)
- **Git**: Usually pre-installed, or `sudo apt install git` (Ubuntu/Debian)
- **Make**: `sudo apt install make` (Ubuntu/Debian)

### Installation Commands

#### Ubuntu/Debian
```bash
# Update package list
sudo apt update

# Install Go
sudo apt install golang-go

# Install Git (if not present)
sudo apt install git

# Install Make
sudo apt install make
```

#### CentOS/RHEL/Fedora
```bash
# Install Go
sudo dnf install golang

# Install Git (if not present)
sudo dnf install git

# Install Make
sudo dnf install make
```

#### macOS
```bash
# Using Homebrew
brew install go
brew install git
brew install make
```

## 🚀 **Running the Bot**

### Quick Start
```bash
# Clone the repository
git clone <your-repo-url>
cd xinli-xlsx-to-ics

# Make the script executable
chmod +x run_bot.sh

# Copy and configure environment
cp config.env.example .env
# Edit .env with your bot token

# Run the bot
./run_bot.sh
```

### Manual Setup
```bash
# Set environment variable
export TELEGRAM_BOT_TOKEN="your_bot_token_here"

# Run the bot
go run .

# Or use .env file
source .env
go run .
```

## 🛠️ **Development Commands**

### Build
```bash
# Build for current platform
go build .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o xlsx-to-ics-linux .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o xlsx-to-ics-macos .

# Build for Windows (if needed)
GOOS=windows GOARCH=amd64 go build -o xlsx-to-ics.exe .
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## 🐳 **Docker Development**

### Build Docker Image
```bash
# Build the image
docker build -t xlsx-to-ics-bot .

# Run locally
docker run -p 8080:8080 --env-file .env xlsx-to-ics-bot

# Run in background
docker run -d -p 8080:8080 --env-file .env --name bot xlsx-to-ics-bot

# Stop and remove
docker stop bot
docker rm bot
```

### Docker Compose (Optional)
Create `docker-compose.yml`:
```yaml
version: '3.8'
services:
  bot:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: unless-stopped
```

Run with:
```bash
docker-compose up -d
docker-compose down
```

## 🔧 **Troubleshooting**

### Permission Issues
```bash
# Fix script permissions
chmod +x run_bot.sh

# Fix file ownership
sudo chown -R $USER:$USER .

# Fix directory permissions
chmod 755 .
```

### Go Issues
```bash
# Clear Go cache
go clean -cache

# Reset Go modules
rm go.sum
go mod tidy

# Check Go version
go version
```

### Environment Issues
```bash
# Check environment variables
env | grep TELEGRAM

# Test .env loading
source .env && echo "Token: $TELEGRAM_BOT_TOKEN"

# Debug environment
set -x
source .env
set +x
```

### Network Issues
```bash
# Test connectivity
curl -I https://api.telegram.org

# Check ports
netstat -tlnp | grep :8080

# Test health endpoint
curl http://localhost:8080/health
```

## 📁 **File Permissions**

### Recommended Permissions
```bash
# Directories
chmod 755 .

# Scripts
chmod 755 run_bot.sh

# Source files
chmod 644 *.go

# Config files
chmod 600 .env
chmod 644 config.env.example

# Documentation
chmod 644 *.md
```

## 🚀 **Production Deployment**

### Systemd Service (Linux)
Create `/etc/systemd/system/xlsx-bot.service`:
```ini
[Unit]
Description=XLSX to ICS Telegram Bot
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/path/to/your/bot
ExecStart=/path/to/your/bot/xlsx-to-ics
Restart=always
RestartSec=10
EnvironmentFile=/path/to/your/bot/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable xlsx-bot
sudo systemctl start xlsx-bot
sudo systemctl status xlsx-bot
```

### Logs
```bash
# View logs
sudo journalctl -u xlsx-bot -f

# View recent logs
sudo journalctl -u xlsx-bot -n 50
```

## 🔒 **Security Best Practices**

### File Permissions
```bash
# Restrict .env access
chmod 600 .env

# Restrict config access
chmod 600 config.env.example

# Secure logs
chmod 640 *.log
```

### Environment Variables
```bash
# Never commit secrets
echo ".env" >> .gitignore

# Use secure random tokens
openssl rand -hex 32

# Validate environment
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "Error: TELEGRAM_BOT_TOKEN not set"
    exit 1
fi
```

## 📚 **Useful Commands**

### Monitoring
```bash
# Watch logs
tail -f bot.log

# Monitor processes
ps aux | grep xlsx-to-ics

# Monitor resources
htop
iotop
```

### Maintenance
```bash
# Clean build artifacts
go clean

# Update dependencies
go get -u ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## 🆘 **Getting Help**

- **Go Documentation**: [golang.org/doc](https://golang.org/doc/)
- **Linux Commands**: [linuxcommand.org](https://linuxcommand.org/)
- **Docker Docs**: [docs.docker.com](https://docs.docker.com/)
- **Systemd**: [systemd.io](https://systemd.io/)
