package responder

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	Code        int
	Payload     any
	ContentType string
}

func Send(w http.ResponseWriter, response *Response) {
	if response.Payload == nil && response.Code != http.StatusNoContent {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("response payload is null"))
		return
	}

	w.Header().Set("Content-Type", response.ContentType)
	w.WriteHeader(response.Code)
	_ = json.NewEncoder(w).Encode(response.Payload)
}

func WrongBodyFormat(response *Response, err error) {
	response.Code = http.StatusBadRequest
	response.Payload = errors.New("wrong body format: " + err.Error())
}

func NotFound(response *Response) {
	response.Code = http.StatusNotFound
	response.Payload = errors.New("not found")
}

func InternalServerError(response *Response, err error) {
	response.Code = http.StatusInternalServerError
	response.Payload = errors.New("internal server error: " + err.Error())
}
