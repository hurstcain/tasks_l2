package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

import "github.com/hurstcain/tasks_l2/develop/dev11/internal/model"

func TestCache_CreateEvent(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-03-22", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)

	validCreate, _ := model.NewEvent("3", "1", "2022-03-22", "abcdee")
	invalidCreate, _ := model.NewEvent("1", "1", "2022-03-22", "abcdee")

	validTestData := []struct {
		create   model.Event
		expected []model.Event
	}{
		{
			create: validCreate,
			expected: []model.Event{
				event1, event2, validCreate,
			},
		},
	}

	invalidTestData := []struct {
		create   model.Event
		expected []model.Event
	}{
		{
			create: invalidCreate,
			expected: []model.Event{
				event1, event2, validCreate,
			},
		},
	}

	for _, data := range validTestData {
		err := cache.CreateEvent(data.create)
		assert.Equal(t, data.expected, cache.events)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		err := cache.CreateEvent(data.create)
		assert.Equal(t, data.expected, cache.events)
		assert.Error(t, err)
	}
}

func TestCache_UpdateEvent(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-03-22", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)

	validUpdate, _ := model.NewEvent("1", "1", "2022-03-22", "abcdee")
	invalidUpdate, _ := model.NewEvent("1111", "1", "2022-03-22", "abcdee")

	validTestData := []struct {
		update   model.Event
		expected []model.Event
	}{
		{
			update: validUpdate,
			expected: []model.Event{
				validUpdate, event2,
			},
		},
	}

	invalidTestData := []struct {
		update   model.Event
		expected []model.Event
	}{
		{
			update: invalidUpdate,
			expected: []model.Event{
				validUpdate, event2,
			},
		},
	}

	for _, data := range validTestData {
		err := cache.UpdateEvent(data.update)
		assert.Equal(t, data.expected, cache.events)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		err := cache.UpdateEvent(data.update)
		assert.Equal(t, data.expected, cache.events)
		assert.Error(t, err)
	}
}

func TestCache_DeleteEvent(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-03-22", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	event3, _ := model.NewEvent("3", "23", "2022-03-28", "1234")
	event4, _ := model.NewEvent("4", "1", "2022-09-09", "1234")
	event5, _ := model.NewEvent("5", "1", "2022-03-26", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)
	cache.CreateEvent(event3)
	cache.CreateEvent(event4)
	cache.CreateEvent(event5)

	validTestData := []struct {
		eventId  string
		expected []model.Event
	}{
		{
			eventId: "2",
			expected: []model.Event{
				event1, event3, event4, event5,
			},
		},
	}

	invalidTestData := []struct {
		eventId  string
		expected []model.Event
	}{
		{
			eventId: "300",
			expected: []model.Event{
				event1, event3, event4, event5,
			},
		},
	}

	for _, data := range validTestData {
		err := cache.DeleteEvent(data.eventId)
		assert.Equal(t, data.expected, cache.events)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		err := cache.DeleteEvent(data.eventId)
		assert.Equal(t, data.expected, cache.events)
		assert.Error(t, err)
	}
}

func TestCache_GetEventsForDay(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-03-22", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	event3, _ := model.NewEvent("3", "23", "2022-03-28", "1234")
	event4, _ := model.NewEvent("4", "1", "2022-09-09", "1234")
	event5, _ := model.NewEvent("5", "1", "2022-03-26", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)
	cache.CreateEvent(event3)
	cache.CreateEvent(event4)
	cache.CreateEvent(event5)

	validTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId: "1",
			date:   "2022-09-09",
			expected: []model.Event{
				event2, event4,
			},
		},
		{
			userId: "23",
			date:   "2022-03-28",
			expected: []model.Event{
				event3,
			},
		},
	}

	invalidTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId:   "10",
			date:     "2022-10-01",
			expected: nil,
		},
		{
			userId:   "1",
			date:     "2022-03-01",
			expected: nil,
		},
	}

	for _, data := range validTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForDay(data.userId, date)
		assert.Equal(t, data.expected, res)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForDay(data.userId, date)
		assert.Nil(t, res)
		assert.Error(t, err)
	}
}

func TestCache_GetEventsForWeek(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-03-22", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	event3, _ := model.NewEvent("3", "23", "2022-03-28", "1234")
	event4, _ := model.NewEvent("4", "1", "2022-10-22", "1234")
	event5, _ := model.NewEvent("5", "1", "2022-03-26", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)
	cache.CreateEvent(event3)
	cache.CreateEvent(event4)
	cache.CreateEvent(event5)

	validTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId: "1",
			date:   "2022-03-21",
			expected: []model.Event{
				event1, event5,
			},
		},
		{
			userId: "23",
			date:   "2022-03-29",
			expected: []model.Event{
				event3,
			},
		},
	}

	invalidTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId:   "10",
			date:     "2022-10-01",
			expected: nil,
		},
		{
			userId:   "1",
			date:     "2022-03-01",
			expected: nil,
		},
	}

	for _, data := range validTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForWeek(data.userId, date)
		assert.Equal(t, data.expected, res)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForWeek(data.userId, date)
		assert.Nil(t, res)
		assert.Error(t, err)
	}
}

func TestCache_GetEventsForMonth(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2022-10-09", "1234")
	event2, _ := model.NewEvent("2", "1", "2022-09-09", "1234")
	event3, _ := model.NewEvent("3", "23", "2022-10-09", "1234")
	event4, _ := model.NewEvent("4", "1", "2022-10-22", "1234")
	event5, _ := model.NewEvent("5", "1", "2022-10-01", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)
	cache.CreateEvent(event3)
	cache.CreateEvent(event4)
	cache.CreateEvent(event5)

	validTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId: "1",
			date:   "2022-10-01",
			expected: []model.Event{
				event1, event4, event5,
			},
		},
		{
			userId: "23",
			date:   "2022-10-20",
			expected: []model.Event{
				event3,
			},
		},
	}

	invalidTestData := []struct {
		userId   string
		date     string
		expected []model.Event
	}{
		{
			userId:   "10",
			date:     "2022-10-01",
			expected: nil,
		},
		{
			userId:   "1",
			date:     "2022-11-20",
			expected: nil,
		},
	}

	for _, data := range validTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForMonth(data.userId, date)
		assert.Equal(t, data.expected, res)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		date, _ := time.Parse(model.DateLayout, data.date)
		res, err := cache.GetEventsForMonth(data.userId, date)
		assert.Nil(t, res)
		assert.Error(t, err)
	}
}

func TestCache_checkEventIdExistence(t *testing.T) {
	event1, _ := model.NewEvent("1", "1", "2020-09-09", "1234")
	event2, _ := model.NewEvent("2", "1", "2020-09-09", "1234")
	event3, _ := model.NewEvent("3", "23", "2020-09-09", "1234")
	event4, _ := model.NewEvent("4", "1", "2020-09-09", "1234")
	event5, _ := model.NewEvent("5", "20", "2020-09-09", "1234")
	cache := NewCache()
	cache.CreateEvent(event1)
	cache.CreateEvent(event2)
	cache.CreateEvent(event3)
	cache.CreateEvent(event4)
	cache.CreateEvent(event5)

	validTestData := []struct {
		id           string
		expectedId   int
		expectedFlag bool
	}{
		{
			id:           "2",
			expectedId:   1,
			expectedFlag: true,
		},
		{
			id:           "20",
			expectedId:   0,
			expectedFlag: false,
		},
	}

	for _, data := range validTestData {
		resId, resFlag := cache.checkEventIdExistence(data.id)
		assert.Equal(t, data.expectedId, resId)
		assert.Equal(t, data.expectedFlag, resFlag)
	}
}

func TestAreEventsEmpty(t *testing.T) {
	validTestData := []struct {
		events   []model.Event
		expected bool
	}{
		{
			events:   make([]model.Event, 0),
			expected: true,
		},
		{
			events:   make([]model.Event, 2),
			expected: false,
		},
	}

	for _, data := range validTestData {
		res := areEventsEmpty(data.events)
		assert.Equal(t, data.expected, res)
	}
}
