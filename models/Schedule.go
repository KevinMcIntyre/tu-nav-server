package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KevinMcIntyre/tu-nav-server/utils"
)

type ScheduleRequest struct {
	StudentID string `json:"sid"`
	TempleUID string `json:"tuid"`
	Password  string `json:"password"`
}

type ScheduleError struct {
	InvalidStudentID bool `json:"sid"`
	InvalidTempleUID bool `json:"tuid"`
	InvalidPassword  bool `json:"password"`
}

type Schedule struct {
	Courses []Course `json:"courses"`
}

type Course struct {
	Name       string `json:"course"`
	Title      string `json:"title"`
	Section    string `json:"section"`
	BuildingID string `json:"buildingId"`
	Room       string `json:"room"`
	Weekday    string `json:"day"`
	StartTime  string `json:"start"`
	EndTime    string `json:"end"`
}

func ConvertTUPortalResponse(response map[string]interface{}) *Schedule {
	for k, v := range response {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
	return nil
}

func (scheduleRequest ScheduleRequest) Validate() *ScheduleError {
	var scheduleErr ScheduleError
	passedValidation := true
	if utils.IsEmptyString(scheduleRequest.StudentID) {
		scheduleErr.InvalidStudentID = true
		passedValidation = false
	}
	if utils.IsEmptyString(scheduleRequest.TempleUID) {
		scheduleErr.InvalidTempleUID = true
		passedValidation = false

	}
	if utils.IsEmptyString(scheduleRequest.Password) {
		scheduleErr.InvalidPassword = true
		passedValidation = false
	}
	if passedValidation {
		return nil
	}
	return &scheduleErr
}

func (scheduleRequest ScheduleRequest) CallTUPortal() (*Schedule, error) {
	client := &http.Client{}
	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://prd-mobile.temple.edu/banner-mobileserver/api/2.0/courses/calendarview/%s", scheduleRequest.StudentID),
		bytes.NewBufferString(""),
	)
	request.SetBasicAuth(scheduleRequest.TempleUID, scheduleRequest.Password)
	request.Header.Add("Content-Type", "application/json")
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		return nil, errors.New("Received a " + response.Status + " response from TUPortal")
	}

	courses := convertTUPortalResponse(response.Body)

	var schedule Schedule
	schedule.Courses = courses

	return &schedule, nil
}

func convertTUPortalResponse(response io.ReadCloser) []Course {
	var schedule map[string][]map[string][]map[string]interface{}
	decoder := json.NewDecoder(response)
	decoder.Decode(&schedule)
	var courses []Course
	for _, courseDay := range schedule["coursesDays"] {
		for _, courseMeeting := range courseDay["coursesMeetings"] {
			var course Course
			course.Name = courseMeeting["courseName"].(string)
			course.Title = courseMeeting["sectionTitle"].(string)
			course.Section = courseMeeting["courseSectionNumber"].(string)
			course.BuildingID = courseMeeting["buildingId"].(string)
			course.Room = courseMeeting["room"].(string)
			startDate, _ := time.Parse("2006-01-02T15:04:05Z", courseMeeting["start"].(string))
			endDate, _ := time.Parse("2006-01-02T15:04:05Z", courseMeeting["end"].(string))
			course.Weekday = startDate.Weekday().String()
			course.StartTime = utils.ConvertTimeObjToTimeString(startDate)
			course.EndTime = utils.ConvertTimeObjToTimeString(endDate)
			courses = append(courses, course)
		}
	}
	return courses
}
