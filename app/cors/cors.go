package cors

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func Apply(router *mux.Router) http.Handler {
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-with"})
	origins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(methods, headers, origins)(router)
}
