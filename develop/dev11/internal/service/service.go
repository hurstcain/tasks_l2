package service

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

import (
	"github.com/hurstcain/tasks_l2/develop/dev11/internal/cache"
	"github.com/hurstcain/tasks_l2/develop/dev11/internal/config"
	"github.com/hurstcain/tasks_l2/develop/dev11/internal/model"
	"github.com/hurstcain/tasks_l2/develop/dev11/internal/service/logger"
)

const (
	ParamEventId      = "event_id"
	ParamUserId       = "user_id"
	ParamDate         = "date"
	ParamEventContent = "event_content"
)

type Service struct {
	server http.Server
	cache  *cache.Cache
	logger *logger.Logger
}

func NewService() *Service {
	addr := config.Ip + ":" + config.Port

	service := &Service{
		server: http.Server{
			Addr: addr,
		},
		cache:  cache.NewCache(),
		logger: logger.NewLogger(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", service.CreateEvent)
	mux.HandleFunc("/update_event", service.UpdateEvent)
	mux.HandleFunc("/delete_event", service.DeleteEvent)
	mux.HandleFunc("/events_for_day", service.GetEventsForDay)
	mux.HandleFunc("/events_for_week", service.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", service.GetEventsForMonth)

	service.server.Handler = service.logger.LogRequest(mux)

	return service
}

func (s *Service) Run() {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return
	}
	if err != nil {
		s.logger.Printf("Error when running server: %s\n", err.Error())
	}
}

func (s *Service) Stop() {
	s.logger.Println("\nClosing server...")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.logger.Printf("Error when closing server: %s\n", err.Error())
	}
}

func (s *Service) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong content-type")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	event, err := parseBodyToEvent(r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Error: %s", err.Error())
		return
	}

	if err := s.cache.CreateEvent(event); err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err = SendPostResponse(w, event)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func (s *Service) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong content-type")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	event, err := parseBodyToEvent(r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Error: %s", err.Error())
		return
	}

	if err := s.cache.UpdateEvent(event); err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err = SendPostResponse(w, event)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func (s *Service) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong content-type")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	eventId := r.PostForm.Get(ParamEventId)
	if err := model.CheckEventId(eventId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Invalid data. Error: %s", err.Error())
		return
	}

	if err := s.cache.DeleteEvent(eventId); err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err := SendDeleteResponse(w)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func (s *Service) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	userId, date, err := parseQueryString(r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Invalid data. Error: %s", err.Error())
		return
	}

	events, err := s.cache.GetEventsForDay(userId, date)
	if err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err = SendGetResponse(w, events)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func (s *Service) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	userId, date, err := parseQueryString(r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Invalid data. Error: %s", err.Error())
		return
	}

	events, err := s.cache.GetEventsForWeek(userId, date)
	if err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err = SendGetResponse(w, events)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func (s *Service) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("Request is not fulfilled. Wrong method")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Request is not fulfilled. Can't parse data from body: %s", err.Error())
		return
	}

	userId, date, err := parseQueryString(r.Form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Printf("Invalid data. Error: %s", err.Error())
		return
	}

	events, err := s.cache.GetEventsForMonth(userId, date)
	if err != nil {
		s.logger.Printf("Business logic error: %s", err.Error())
		responseErr := SendErrorResponse503(w, err)
		if responseErr != nil {
			s.logger.Printf("Error: %s", err.Error())
		}
		return
	}

	err = SendGetResponse(w, events)
	if err != nil {
		s.logger.Printf("Error: %s", err.Error())
	}
}

func parseBodyToEvent(s url.Values) (model.Event, error) {
	eventId := s.Get(ParamEventId)
	userId := s.Get(ParamUserId)
	date := s.Get(ParamDate)
	eventContent := s.Get(ParamEventContent)

	event, err := model.NewEvent(eventId, userId, date, eventContent)
	return event, err
}

func parseQueryString(s url.Values) (string, time.Time, error) {
	userId := s.Get(ParamUserId)
	date := s.Get(ParamDate)

	if err := model.CheckUserId(userId); err != nil {
		return "", time.Time{}, err
	}

	t, err := model.CheckDate(date)
	return userId, t, err
}
