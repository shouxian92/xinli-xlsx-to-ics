#!/bin/bash

echo "Setting up Telegram Bot..."
echo

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "Error: .env file not found!"
    echo "Please copy config.env.example to .env and add your bot token."
    echo
    read -p "Press Enter to exit"
    exit 1
fi

# Load environment variables from .env file
export $(cat .env | grep -v '^#' | xargs)

# Check if token is set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "Error: TELEGRAM_BOT_TOKEN not found in .env file!"
    echo "Please check your .env file configuration."
    echo
    read -p "Press Enter to exit"
    exit 1
fi

echo "Bot token loaded successfully!"
echo "Starting bot..."
echo
echo "Press Ctrl+C to stop the bot"
echo

# Run the bot
go run .
