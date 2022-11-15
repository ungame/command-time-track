package httpext

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"
	MimeJson          = "application/json"
)

type Port int

func (p Port) Addr() string {
	return fmt.Sprintf(":%d", p)
}

type ErrorOutput struct {
	Err string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set(HeaderContentType, MimeJson)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error on encode json:", err.Error())
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	out := ErrorOutput{"unknown error"}
	if err != nil {
		out.Err = err.Error()
	}
	WriteJson(w, status, out)
}
