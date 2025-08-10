package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
	tele "gopkg.in/telebot.v3"
)

const TABLE_ROW_SIZE int = 21

func main() {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	// Bot settings
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// Initialize bot
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle /start command
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Welcome! I can convert Excel timetables to ICS calendar files. Just send me an Excel file (.xlsx) and I'll convert it for you.")
	})

	// Handle /help command
	bot.Handle("/help", func(c tele.Context) error {
		helpText := `Available commands:
/start - Start the bot
/help - Show this help message

To convert an Excel timetable:
1. Send me an Excel file (.xlsx)
2. I'll process it and send back an ICS calendar file
3. You can then import the ICS file into your calendar app

Note: The Excel file should follow the expected format with modules at the top and weekly timetables below.`
		return c.Send(helpText)
	})

	// Handle file uploads
	bot.Handle(tele.OnDocument, func(c tele.Context) error {
		doc := c.Message().Document

		// Check if it's an Excel file
		if !strings.HasSuffix(strings.ToLower(doc.FileName), ".xlsx") {
			return c.Send("Please send an Excel file (.xlsx) for conversion.")
		}

		// Send processing message
		if err := c.Send("Processing your Excel file... Please wait."); err != nil {
			return c.Send("Error: Could not send processing message.")
		}

		// Create temporary file for download
		tempFile, err := os.CreateTemp("", "xlsx_*.xlsx")
		if err != nil {
			return c.Send("Error: Could not create temporary file.")
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Download the file
		if err := bot.Download(&doc.File, tempFile.Name()); err != nil {
			return c.Send("Error: Could not download the file.")
		}

		// Process the Excel file
		icsData, filename, err := processExcelFile(tempFile.Name())
		if err != nil {
			return c.Send(fmt.Sprintf("Error processing file: %v", err))
		}

		// Send the ICS file
		icsDoc := &tele.Document{
			File:     tele.FromReader(strings.NewReader(icsData)),
			FileName: filename + ".ics",
		}

		return c.Send(icsDoc)
	})

	// Handle text messages
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Please send me an Excel file (.xlsx) to convert to ICS format, or use /help for more information.")
	})

	// Start HTTP server for health checks (required by Render)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("HTTP server starting on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Bot started. Press Ctrl+C to stop.")
	bot.Start()
}

// processExcelFile processes the uploaded Excel file and returns ICS data
func processExcelFile(filePath string) (string, string, error) {
	// Open the Excel file
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open Excel file: %w", err)
	}

	if len(wb.Sheets) == 0 {
		return "", "", fmt.Errorf("no sheets found in Excel file")
	}

	rows := wb.Sheets[0].Rows
	if len(rows) < 2 {
		return "", "", fmt.Errorf("insufficient data in Excel file")
	}

	// Generate output filename
	output_filename := rows[0].Cells[0].Value + " - " + rows[1].Cells[0].Value
	if output_filename == " - " {
		output_filename = "Timetable"
	}

	// Build modules and week timetables
	modules, stopIndex := buildModules(rows)
	weeks := buildWeekTimetables(rows, stopIndex)

	// Create ICS calendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetXWRCalName(output_filename)

	// Add events to calendar
	for _, w := range weeks {
		for _, days := range w.lessons {
			for _, d := range days {
				event := cal.AddEvent("xinli_s93@hotmail.com" + uuid.NewString())
				event.SetCreatedTime(time.Now())
				event.SetDtStampTime(time.Now())
				event.SetModifiedAt(time.Now())
				event.SetStartAt(d.startTime)
				event.SetEndAt(d.endTime)

				mod, exists := modules[d.module]
				if exists {
					event.SetSummary(fmt.Sprintf("[%s] %s", mod.code, mod.name))
				} else {
					event.SetSummary(d.module)
				}

				if len(d.location) > 0 {
					event.SetLocation(d.location)
				}
				event.SetDescription(d.description)
				event.SetOrganizer("xinli_s93@hotmail.com")
			}
		}
	}

	// Serialize calendar to ICS format
	icsData := cal.Serialize()
	return icsData, output_filename, nil
}
