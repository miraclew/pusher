package api

import (
	"encoding/json"
	"net/http"
)

func respondOK(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	res.WriteHeader(http.StatusOK)
	response, err := json.Marshal(data)
	if err != nil {
		response, _ = json.Marshal(map[string]interface{}{
			"message": err.Error(),
		})
		res.Write(response)
		return
	}

	res.Write(response)
}

func respondFail(res http.ResponseWriter, statusCode int, message string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
	res.Write(response)
}
