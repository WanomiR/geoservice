package rr

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ReadResponder interface {
	ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
	WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
	WriteJSONError(w http.ResponseWriter, err error, status ...int) error
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ReadRespond struct {
	maxBytes int
}

type ReadRespondOption func(*ReadRespond)

func WithMaxBytes(maxBytes int) ReadRespondOption {
	return func(r *ReadRespond) {
		r.maxBytes = maxBytes
	}
}

func NewReadRespond(options ...ReadRespondOption) *ReadRespond {
	rr := &ReadRespond{}

	for _, option := range options {
		option(rr)
	}
	return rr
}

func (rr *ReadRespond) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	if rr.maxBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, int64(rr.maxBytes))
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

func (rr *ReadRespond) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (rr *ReadRespond) WriteJSONError(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	response := &JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return rr.WriteJSON(w, statusCode, response)
}
