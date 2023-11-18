package middlewares

import (
	"log"
	"mime"
	"net/http"

	"cards/internal/models"
	writeJSONresponse "cards/internal/pkg/writeJSONresponse"
)

func EnforceJSONHandler(next http.Handler) http.Handler {
	log.Print("Executing enforceJSONHandler")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				writeJSONresponse.WriteJSONresponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil, models.Pagination{}, "Malformed Content-Type header")
				return
			}

			if mt != "application/json" {
				writeJSONresponse.WriteJSONresponse(w, http.StatusUnsupportedMediaType, http.StatusText(http.StatusBadRequest), nil, models.Pagination{}, "Content-Type header must be application/json")
				return
			}
		}

		next.ServeHTTP(w, r)
		log.Print("Executing enforceJSONHandler again")
	})
}