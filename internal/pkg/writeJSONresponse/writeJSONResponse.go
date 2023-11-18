package pkg

import (
	"cards/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func WriteJSONresponse(r http.ResponseWriter, status int, message string, result any, p models.Pagination, err string) error {
	log.Print("Executing writeJSONresponse")

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)

	response := make(map[string]interface{})
	response["status"] = strconv.Itoa(status)
	response["message"] = message
	response["result"] = result
	response["error"] = err
	response["pagination"] = p
	
	return json.NewEncoder(r).Encode(response)
}