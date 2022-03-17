package service

import (
	"encoding/json"
	"net/http"
)

import "github.com/hurstcain/tasks_l2/develop/dev11/internal/model"

type PostResponse struct {
	Result model.Event `json:"result"`
}

type DeleteResponse struct {
	Result string `json:"result"`
}

type GetResponse struct {
	Result []model.Event `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendPostResponse(w http.ResponseWriter, event model.Event) error {
	resultResponse := PostResponse{
		Result: event,
	}

	response, err := json.Marshal(resultResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func SendDeleteResponse(w http.ResponseWriter) error {
	resultResponse := DeleteResponse{
		Result: "Event has been deleted",
	}

	response, err := json.Marshal(resultResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func SendGetResponse(w http.ResponseWriter, events []model.Event) error {
	resultResponse := GetResponse{
		Result: events,
	}

	response, err := json.Marshal(resultResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func SendErrorResponse503(w http.ResponseWriter, err error) error {
	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	response, err := json.Marshal(errorResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.Error(w, string(response), http.StatusServiceUnavailable)

	return nil
}
