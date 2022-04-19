package middlewares

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if jsonResponse, err := json.Marshal(response); err != nil {
		panic(err)
	} else {
		w.Write(jsonResponse)
	}
}