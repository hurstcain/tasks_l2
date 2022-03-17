package model

import (
	"fmt"
	"time"
)

const DateLayout = "2006-01-02"

type Event struct {
	EventId      string    `json:"event_id"`
	UserId       string    `json:"user_id"`
	DateString   string    `json:"date"`
	Date         time.Time `json:"-"`
	EventContent string    `json:"event_content"`
}

func NewEvent(eventId, userId, dateString, eventContent string) (Event, error) {
	if err := CheckEventId(eventId); err != nil {
		return Event{}, err
	}

	if err := CheckUserId(userId); err != nil {
		return Event{}, err
	}

	date, err := CheckDate(dateString)
	if err != nil {
		return Event{}, err
	}

	return Event{
		EventId:      eventId,
		UserId:       userId,
		DateString:   dateString,
		Date:         date,
		EventContent: eventContent,
	}, nil
}

func CheckEventId(eventId string) error {
	if isEmpty(eventId) {
		return fmt.Errorf("EventId is empty")
	}

	return nil
}

func CheckUserId(userId string) error {
	if isEmpty(userId) {
		return fmt.Errorf("UserId is empty")
	}

	return nil
}

func CheckDate(date string) (time.Time, error) {
	if isEmpty(date) {
		return time.Time{}, fmt.Errorf("date is empty")
	}

	t, err := time.Parse(DateLayout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("wrong date format: %s. Valid date format: YYYY-MM-DD", err.Error())
	}

	return t, nil
}

func isEmpty(s string) bool {
	return s == ""
}
