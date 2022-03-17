package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	validTestData := []struct {
		eventId      string
		userId       string
		dateString   string
		eventContent string
	}{
		{
			eventId:      "1",
			userId:       "1",
			dateString:   "2022-07-08",
			eventContent: "1",
		},
	}

	invalidTestData := []struct {
		eventId      string
		userId       string
		dateString   string
		eventContent string
	}{
		{
			eventId:      "",
			userId:       "",
			dateString:   "",
			eventContent: "1",
		},
		{
			eventId:      "",
			userId:       "1",
			dateString:   "2022-03-04",
			eventContent: "1",
		},
		{
			eventId:      "1",
			userId:       "",
			dateString:   "2022-03-04",
			eventContent: "1",
		},
		{
			eventId:      "1",
			userId:       "1",
			dateString:   "2022-033-04",
			eventContent: "1",
		},
	}

	for _, data := range validTestData {
		res, err := NewEvent(data.eventId, data.userId, data.dateString, data.eventContent)
		date, _ := CheckDate(data.dateString)
		eventExpected := Event{
			EventId:      data.eventId,
			UserId:       data.userId,
			DateString:   data.dateString,
			Date:         date,
			EventContent: data.eventContent,
		}
		assert.Equal(t, eventExpected, res)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		_, err := NewEvent(data.eventId, data.userId, data.dateString, data.eventContent)
		assert.Error(t, err)
	}
}

func TestCheckEventId(t *testing.T) {
	validTestData := []struct {
		id string
	}{
		{
			id: "1234",
		},
	}

	invalidTestData := []struct {
		id string
	}{
		{
			id: "",
		},
	}

	for _, data := range validTestData {
		err := CheckEventId(data.id)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		err := CheckEventId(data.id)
		assert.Error(t, err)
	}
}

func TestCheckUserId(t *testing.T) {
	validTestData := []struct {
		id string
	}{
		{
			id: "1234",
		},
	}

	invalidTestData := []struct {
		id string
	}{
		{
			id: "",
		},
	}

	for _, data := range validTestData {
		err := CheckUserId(data.id)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		err := CheckUserId(data.id)
		assert.Error(t, err)
	}
}

func TestCheckDate(t *testing.T) {
	validTestData := []struct {
		date     string
		expected time.Time
	}{
		{
			date: "2020-09-10",
		},
		{
			date: "2023-07-13",
		},
	}

	invalidTestData := []struct {
		date string
	}{
		{
			date: "",
		},
		{
			date: "2022-99-01",
		},
		{
			date: "2022-01-90",
		},
		{
			date: "abc",
		},
	}

	for _, data := range validTestData {
		resDate, err := CheckDate(data.date)
		date, _ := time.Parse(DateLayout, data.date)
		assert.NoError(t, err)
		assert.Equal(t, date, resDate)
	}

	for _, data := range invalidTestData {
		_, err := CheckDate(data.date)
		assert.Error(t, err)
	}
}

func TestIsEmpty(t *testing.T) {
	validTestData := []struct {
		s        string
		expected bool
	}{
		{
			s:        "",
			expected: true,
		},
		{
			s:        "abcd",
			expected: false,
		},
	}

	for _, data := range validTestData {
		res := isEmpty(data.s)
		assert.Equal(t, data.expected, res)
	}
}
