package main

import (
	"net/http"
	"time"
)

func (app *Config) messageReceived(message Data) error {
	record := MotorSpeedRecord{
		Time:       time.Now(),
		MotorSpeed: message.MotorSpeed,
	}
	return app.InsertMotorSpeedRecord(record)
}

type RequestPayload struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (app *Config) FetchData(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
	}

	layout := "2006-01-02T15:04:05.000Z"
	parsedStart, err := time.Parse(layout, requestPayload.Start)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	parsedEnd, err := time.Parse(layout, requestPayload.End)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	records, err := app.FetchMotorSpeedRecords(parsedStart, parsedEnd)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	responsePayload := jsonResponse{
		Message: "fetched",
		Data:    records,
	}

	app.writeJson(w, http.StatusOK, responsePayload)
}
