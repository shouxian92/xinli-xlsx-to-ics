# Quick Setup Guide

## 1. Get a Telegram Bot Token

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the token provided by BotFather

## 2. Configure the Bot

1. Copy `config.env.example` to `.env`
2. Edit `.env` and add your bot token:
   ```
   TELEGRAM_BOT_TOKEN=your_actual_bot_token_here
   ```

## 3. Run the Bot

### Option 1: Using Make (Recommended for Linux/macOS)
```bash
# Setup development environment
make setup

# Run the bot
make run
```

### Option 2: Shell Script
```bash
# Make script executable (first time only)
chmod +x run_bot.sh

# Run the bot
./run_bot.sh
```

### Option 3: Manual
```bash
# Set environment variable
export TELEGRAM_BOT_TOKEN="your_token_here"

# Run the bot
go run .

# Or use .env file
source .env
go run .
```

## 4. Test the Bot

1. Find your bot on Telegram (using the username you created)
2. Send `/start` to begin
3. Upload an Excel file (.xlsx) to test the conversion
4. The bot will send back an ICS file

## Troubleshooting

- **"TELEGRAM_BOT_TOKEN environment variable is required"**: Make sure you've created the `.env` file with your token
- **Bot not responding**: Ensure you've started a conversation with the bot first
- **File processing errors**: Check that your Excel file follows the expected format

## Expected Excel Format

The bot expects Excel files with:
- Module table at the top (codes, names, credits, instructors)
- Weekly timetables below (starting with "TIME" headers)
- 21 rows per week timetable
- Columns representing different days of the week
