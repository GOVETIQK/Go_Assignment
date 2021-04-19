package util

import (
	"encoding/json"
	"net/http"

	"github.com/go-errors/errors"
)

type Error struct {
	Msg   string `json:"errorMsg"`
	Trace string `json:"trace"`
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	emsg := Error{
		Msg:   err.Error(),
		Trace: err.(*errors.Error).ErrorStack(),
	}
	RespondWithJSON(w, code, emsg)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
