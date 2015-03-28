package api

import (
	"encoding/json"
	"net/http"
)

func respondOk(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	response, err := json.Marshal(data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		response, _ = json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		res.Write(response)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func respondFail(res http.ResponseWriter, statusCode int, message string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	response, _ := json.Marshal(map[string]interface{}{
		"message": message,
	})
	res.Write(response)
}
