package middlewares

import (
	"log"
	"net/http"
	"time"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusCodeRecorder) WriteHeader(status int) {
	r.WriteHeader(status)
	r.status = status
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			start    = time.Now()
			proto    = request.Proto
			method   = request.Method
			uri      = request.RequestURI
			recorder = &statusCodeRecorder{writer, 0}
		)

		next(writer, request)

		log.Printf("%s %s %s %d %s", proto, method, uri, recorder.status, time.Since(start).String())
	}
}
