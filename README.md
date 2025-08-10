# XLSX to ICS Telegram Bot

A Telegram bot that converts Excel timetable files (.xlsx) to ICS calendar format (.ics) files. Users can upload Excel files directly to the bot and receive converted calendar files in return.

## Features

- **File Upload Handling**: Accepts .xlsx files from users
- **Dynamic Processing**: Processes any Excel file following the expected format
- **ICS Generation**: Converts Excel timetables to standard calendar format
- **User-Friendly**: Simple commands and helpful error messages
- **Real-time Processing**: Converts files on-demand

## Prerequisites

- Go 1.22.0 or higher
- A Telegram Bot Token (obtained from [@BotFather](https://t.me/botfather))

## Installation

1. **Clone the repository**:
   ```bash
   git clone <your-repo-url>
   cd xinli-xlsx-to-ics
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up your bot token**:
   - Copy `config.env.example` to `.env`
   - Edit `.env` and add your Telegram bot token:
     ```
     TELEGRAM_BOT_TOKEN=your_actual_bot_token_here
     ```

## Getting a Telegram Bot Token

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the token provided by BotFather
5. Add it to your `.env` file

## Usage

### Running the Bot

#### Option 1: Using Make (Recommended for Linux/macOS)
```bash
# Setup development environment
make setup

# Run the bot
make run

# Or build and run manually
make build
./xlsx-to-ics
```

#### Option 2: Using the provided script
```bash
# Make script executable (first time only)
chmod +x run_bot.sh

# Run the bot
./run_bot.sh
```

#### Option 3: Manual setup
```bash
# Set environment variable and run
export TELEGRAM_BOT_TOKEN="your_token_here"
go run .

# Or use the .env file
source .env
go run .
```

#### Option 4: Docker
```bash
# Build and run with Docker
make docker-build
make docker-run
```

### Bot Commands

- `/start` - Welcome message and instructions
- `/help` - Detailed help information

### How to Use

1. **Start the bot** with `/start`
2. **Upload an Excel file** (.xlsx) to the bot
3. **Wait for processing** - the bot will show a "Processing..." message
4. **Receive the ICS file** - the bot will send back the converted calendar file
5. **Import to your calendar** - download and import the .ics file to your preferred calendar app

## Expected Excel Format

The bot expects Excel files with the following structure:

1. **Module Table** (at the top):
   - Row 1: Module codes
   - Row 2: Module names
   - Row 3: Credits
   - Row 4: Lead instructors
   - Row 5+: Additional module information

2. **Weekly Timetables** (below modules):
   - Each week starts with a "TIME" header
   - 21 rows per week timetable
   - Columns represent different days of the week
   - Time slots are filled with module codes and details

## File Structure

```
xinli-xlsx-to-ics/
├── main.go              # Main bot logic and handlers
├── structs.go           # Data structures for modules and timetables
├── xlsx.go              # Excel processing functions
├── go.mod               # Go module dependencies
├── config.env.example   # Example configuration
├── run_bot.sh           # Linux/macOS startup script
├── Makefile             # Development commands and automation
├── Dockerfile           # Docker configuration for deployment
├── render.yaml          # Render deployment configuration
├── render-simple.yaml   # Alternative Go-native deployment
├── LINUX_DEVELOPMENT.md # Linux-specific development guide
├── DEPLOYMENT.md        # Render deployment guide
├── SETUP.md             # Quick setup guide
└── README.md            # This file
```

## Dependencies

- `gopkg.in/telebot.v3` - Telegram bot framework
- `github.com/tealeg/xlsx` - Excel file processing
- `github.com/arran4/golang-ical` - ICS calendar generation
- `github.com/google/uuid` - Unique identifier generation

## Error Handling

The bot includes comprehensive error handling for:
- Invalid file types (non-Excel files)
- File download failures
- Excel processing errors
- Insufficient data in files
- Missing sheets

## Security Notes

- Never commit your `.env` file with real tokens
- The bot only accepts .xlsx files
- Files are processed in memory and not stored permanently
- Each user's files are processed independently

## Troubleshooting

### Common Issues

1. **"TELEGRAM_BOT_TOKEN environment variable is required"**
   - Make sure you've set the environment variable or created a .env file

2. **"Could not download the file"**
   - Check your internet connection
   - Ensure the file isn't too large (Telegram has file size limits)

3. **"Error processing file"**
   - Verify the Excel file follows the expected format
   - Check that the file isn't corrupted

4. **Bot not responding**
   - Ensure the bot is running
   - Check that you've started a conversation with the bot first

## Development

### Prerequisites

- Go 1.22.0 or higher
- Telegram Bot Token from [@BotFather](https://t.me/botfather)
- Linux/macOS environment (Windows users can use WSL)

### Quick Start for Linux/macOS

```bash
# Clone and setup
git clone <your-repo-url>
cd xinli-xlsx-to-ics

# Setup development environment
make setup

# Edit .env with your bot token
nano .env

# Run the bot
make run
```

### Available Make Commands

```bash
make help          # Show all available commands
make setup         # Setup development environment
make build         # Build the bot binary
make run           # Run the bot locally
make test          # Run tests
make clean         # Clean build artifacts
make deps          # Download dependencies
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make deploy        # Build for deployment
make check         # Run all quality checks
make status        # Show bot status
```

### Modifying the Bot

1. **Add new commands**: Add new handlers in `main.go`
2. **Modify Excel processing**: Update functions in `xlsx.go`
3. **Change data structures**: Modify `structs.go`
4. **Add new features**: Extend the bot handlers as needed

For detailed Linux development information, see [LINUX_DEVELOPMENT.md](LINUX_DEVELOPMENT.md).

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
