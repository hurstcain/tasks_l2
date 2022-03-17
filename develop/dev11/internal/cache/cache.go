package cache

import (
	"fmt"
	"sync"
	"time"
)

import "github.com/hurstcain/tasks_l2/develop/dev11/internal/model"

type Cache struct {
	events []model.Event
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		events: make([]model.Event, 0),
	}
}

func (c *Cache) CreateEvent(event model.Event) error {
	if _, exists := c.checkEventIdExistence(event.EventId); exists {
		return fmt.Errorf("event with this id already exists")
	}

	c.Lock()
	c.events = append(c.events, event)
	c.Unlock()

	return nil
}

func (c *Cache) UpdateEvent(event model.Event) error {
	i, exists := c.checkEventIdExistence(event.EventId)
	if !exists {
		return fmt.Errorf("event with this id doesn't exist")
	}

	c.Lock()
	c.events[i] = event
	c.Unlock()

	return nil
}

func (c *Cache) DeleteEvent(eventId string) error {
	i, exists := c.checkEventIdExistence(eventId)
	if !exists {
		return fmt.Errorf("event with this id doesn't exist")
	}

	c.Lock()
	c.events = append(c.events[:i], c.events[i+1:]...)
	c.Unlock()

	return nil
}

func (c *Cache) GetEventsForDay(userId string, date time.Time) ([]model.Event, error) {
	eventsForDay := make([]model.Event, 0)

	c.RLock()
	for _, event := range c.events {
		if event.UserId == userId && event.Date == date {
			eventsForDay = append(eventsForDay, event)
		}
	}
	c.RUnlock()

	if areEventsEmpty(eventsForDay) {
		return nil, fmt.Errorf("no events for this day")
	}

	return eventsForDay, nil
}

func (c *Cache) GetEventsForWeek(userId string, date time.Time) ([]model.Event, error) {
	eventsForWeek := make([]model.Event, 0)
	year, week := date.ISOWeek()

	c.RLock()
	for _, event := range c.events {
		eventYear, eventWeek := event.Date.ISOWeek()
		if event.UserId == userId && eventYear == year && eventWeek == week {
			eventsForWeek = append(eventsForWeek, event)
		}
	}
	c.RUnlock()

	if areEventsEmpty(eventsForWeek) {
		return nil, fmt.Errorf("no events for this week")
	}

	return eventsForWeek, nil
}

func (c *Cache) GetEventsForMonth(userId string, date time.Time) ([]model.Event, error) {
	eventsForYear := make([]model.Event, 0)
	year := date.Year()
	month := date.Month()

	c.RLock()
	for _, event := range c.events {
		eventYear := event.Date.Year()
		eventMonth := event.Date.Month()
		if event.UserId == userId && eventYear == year && eventMonth == month {
			eventsForYear = append(eventsForYear, event)
		}
	}
	c.RUnlock()

	if areEventsEmpty(eventsForYear) {
		return nil, fmt.Errorf("no events for this year")
	}

	return eventsForYear, nil
}

func (c *Cache) checkEventIdExistence(eventId string) (int, bool) {
	c.RLock()
	defer c.RUnlock()
	for i, event := range c.events {
		if event.EventId == eventId {
			return i, true
		}
	}

	return 0, false
}

func areEventsEmpty(slice []model.Event) bool {
	return len(slice) == 0
}
