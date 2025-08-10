package main

import "time"

type person string

type module struct {
	code      string
	name      string
	credits   string
	lead      person
	cellColor string
}

type weekTimetable struct {
	startDate time.Time
	endDate   time.Time
	lessons   [][]lesson
}

type lesson struct {
	module      string //module
	description string
	startTime   time.Time
	endTime     time.Time
	location    string
}
