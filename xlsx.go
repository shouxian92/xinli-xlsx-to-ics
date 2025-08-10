package main

import (
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

const TABLE_ROW_SIZE int = 21

// map the modules by module code
func buildModules(rows []*xlsx.Row) (map[string]module, int) {
	modules := map[string]module{}
	// the module table is always at the top of the document, so we just loop until we see the "CODE" header and stop
	// this is a reasonable approach since there aren't many rows before the module table
	current := 0
	for i := range rows {
		if rows[i].Cells[0].Value == "CODE" {
			current = i + 1
			break
		}
	}

	// from this point onwards we iterate till we get an empty row
	for {
		if len(rows[current].Cells) == 0 || rows[current].Cells[0].Value == "" {
			break
		}
		code := rows[current].Cells[0].Value
		name := rows[current].Cells[1].Value
		credits := rows[current].Cells[2].Value
		lead := person(rows[current].Cells[3].Value)
		cellColor := rows[current].Cells[0].GetStyle().Fill.FgColor
		module := module{
			code:      code,
			name:      name,
			credits:   credits,
			lead:      lead,
			cellColor: cellColor,
		}
		modules[code] = module
		current++
	}

	return modules, current
}

// buildWeekTimetables builds the week timetables from the rows of the spreadsheet
// it takes the start index of the week timetable and iterates until it reaches the "TIME" header
// it then builds a week timetable for each week
// it returns a slice of week timetables
func buildWeekTimetables(rows []*xlsx.Row, startIndex int) []weekTimetable {
	weeks := []weekTimetable{}

	for i := startIndex; i < len(rows); i++ {
		// iterate until we reach TIME
		if len(rows[i].Cells) == 0 || rows[i].Cells[0].Value != "Time" {
			continue
		}

		// handle cases in the spreadsheet where the row before "TIME" is not prefixed with "Week" or "EXAM"
		// example of these edge cases are "Recess Week" or "Reading Week"
		// the timetable is empty in these cases, so we skip to the next valid week
		if isEdgeCaseTable(rows, i) {
			nextIndex := findNextValidTableIndex(rows, i)
			if nextIndex != -1 {
				i = nextIndex - 1 // -1 because the loop will increment i
				continue
			} else {
				// If no valid next index found, break out of the loop
				break
			}
		}

		// a table consists of 21 rows and 11 columns, so everytime we build a week
		// we skip 21 rows and continue iterating
		weeks = append(weeks, buildWeek(rows, i))
		i += TABLE_ROW_SIZE
	}
	return weeks
}

// isEdgeCaseTable checks if the current table is an edge case (not a standard Week/EXAM table)
// Edge cases include "Recess Week", "Reading Week", etc. which have empty timetables
func isEdgeCaseTable(rows []*xlsx.Row, timeHeaderRow int) bool {
	if timeHeaderRow == 0 {
		return false
	}

	prevRow := rows[timeHeaderRow-1]
	if len(prevRow.Cells) == 0 {
		return false
	}

	prevRowValue := prevRow.Cells[0].Value
	return !strings.HasPrefix(prevRowValue, "Week") && prevRowValue != "EXAM"
}

// findNextValidTableIndex searches forward from the current position to find the next valid table
// A valid table is one where the row before "TIME" is prefixed with "Week" or "EXAM"
func findNextValidTableIndex(rows []*xlsx.Row, currentTimeHeaderRow int) int {
	for searchRowIndex := currentTimeHeaderRow + 1; searchRowIndex < len(rows); searchRowIndex++ {
		// Check if we found the next valid week/exam
		if isValidWeekOrExamRow(rows, searchRowIndex) {
			return searchRowIndex
		}

		// Check if we found another "TIME" header (indicating another table)
		if isTimeHeader(rows, searchRowIndex) {
			// If we find another TIME header, check if the row before it is valid
			if searchRowIndex > 0 && isValidWeekOrExamRow(rows, searchRowIndex-1) {
				return searchRowIndex
			}
		}
	}

	return -1 // No valid table found
}

// isValidWeekOrExamRow checks if a row represents a valid week or exam (has Week prefix or is EXAM)
func isValidWeekOrExamRow(rows []*xlsx.Row, rowIndex int) bool {
	if rowIndex >= len(rows) || len(rows[rowIndex].Cells) == 0 {
		return false
	}

	rowValue := rows[rowIndex].Cells[0].Value
	return strings.HasPrefix(rowValue, "Week") || rowValue == "EXAM"
}

// isTimeHeader checks if a row is a "Time" header row
func isTimeHeader(rows []*xlsx.Row, rowIndex int) bool {
	if rowIndex >= len(rows) || len(rows[rowIndex].Cells) == 0 {
		return false
	}

	return rows[rowIndex].Cells[0].Value == "Time"
}

func buildWeek(rows []*xlsx.Row, weekStartRow int) weekTimetable {
	// the first column consists of all the time and the next 5 columns are days of the week
	// each day of the week occupies 2 columns
	// Use UTC+8 (GMT+8) directly instead of loading Asia/Singapore timezone
	loc := time.FixedZone("GMT+8", 8*60*60)

	// Add safety check for date parsing
	if len(rows[weekStartRow].Cells) < 2 {
		// Return an empty week timetable with current time as fallback
		return weekTimetable{
			startDate: time.Now().In(loc),
			endDate:   time.Now().In(loc).AddDate(0, 0, 5),
			lessons:   [][]lesson{},
		}
	}

	startDate, err := rows[weekStartRow].Cells[1].GetTime(false)
	if err != nil {
		// Return an empty week timetable with current time as fallback
		return weekTimetable{
			startDate: time.Now().In(loc),
			endDate:   time.Now().In(loc).AddDate(0, 0, 5),
			lessons:   [][]lesson{},
		}
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), loc)

	weekTimetable := weekTimetable{
		startDate: startDate,
		endDate:   startDate.AddDate(0, 0, 5),
		lessons:   [][]lesson{},
	}
	currentRow := weekStartRow + 2

	for column := 1; column < 11; column += 2 {
		lessons := []lesson{}
		// at each day, search for the start of a lesson
		// start of a lesson is defined by the cell color
		for startLessonRow := currentRow; startLessonRow < weekStartRow+TABLE_ROW_SIZE; startLessonRow++ {
			if len(rows[startLessonRow].Cells) == 0 || rows[startLessonRow].Cells[column].GetStyle().Fill.FgColor == "" {
				// nothing to map so we just skip
				continue
			}

			lesson, lessonEndIndex := buildLessonTimeslot(rows, weekStartRow, startLessonRow, column)
			lessons = append(lessons, lesson)
			startLessonRow = lessonEndIndex
		}
		weekTimetable.lessons = append(weekTimetable.lessons, lessons)
	}

	return weekTimetable
}

func buildLessonTimeslot(rows []*xlsx.Row, weekStartRow, lessonStartRow, dayColumn int) (lesson, int) {

	// Add safety checks for array bounds
	if lessonStartRow >= len(rows) || dayColumn >= len(rows[lessonStartRow].Cells) {
		return lesson{}, lessonStartRow
	}

	if lessonStartRow+1 >= len(rows) || dayColumn+1 >= len(rows[lessonStartRow+1].Cells) {
		return lesson{}, lessonStartRow
	}

	// first row always has the module code
	// 1. there is an edge case we have to handle here when there are whitespaces in the module code. replace all whitespaces with an empty string
	code := strings.ReplaceAll(rows[lessonStartRow].Cells[dayColumn].Value, " ", "")
	description := rows[lessonStartRow+1].Cells[dayColumn].Value

	// Use UTC+8 (GMT+8) directly instead of loading Asia/Singapore timezone
	loc := time.FixedZone("GMT+8", 8*60*60)

	// Add safety check for date parsing
	if len(rows[weekStartRow].Cells) <= dayColumn {
		return lesson{}, lessonStartRow
	}

	startDate, err := rows[weekStartRow].Cells[dayColumn].GetTime(false)
	if err != nil {
		return lesson{}, lessonStartRow
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), loc)

	startTime := startDate.Add(8*time.Hour + 30*time.Minute)

	// now, how far away are we from the first row?
	startTime = startTime.Add(time.Duration((lessonStartRow-weekStartRow-2)*30) * time.Minute)
	endTime := startTime.Add(30 * time.Minute)

	lessonCellColor := rows[lessonStartRow].Cells[dayColumn].GetStyle().Fill.FgColor
	endOfLessonIndex := lessonStartRow
	for ; endOfLessonIndex < weekStartRow+TABLE_ROW_SIZE && endOfLessonIndex < len(rows); endOfLessonIndex++ {
		// Keep going down the rows
		if len(rows[endOfLessonIndex].Cells) <= dayColumn || rows[endOfLessonIndex].Cells[dayColumn].GetStyle().Fill.FgColor != lessonCellColor {
			break
		}
		endTime = startTime.Add(time.Duration((endOfLessonIndex-lessonStartRow+1)*30) * time.Minute)
	}

	location := ""
	if len(rows[lessonStartRow+1].Cells) > dayColumn+1 {
		location = rows[lessonStartRow+1].Cells[dayColumn+1].Value
	}

	l := lesson{
		module:      code,
		description: description,
		startTime:   startTime,
		endTime:     endTime,
		location:    location,
	}

	return l, endOfLessonIndex - 1
}
