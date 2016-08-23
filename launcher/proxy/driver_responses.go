package driverhub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const contentType = "application/json; charset=utf-8"

func success(w http.ResponseWriter, value interface{}) {
	body, err := json.Marshal(map[string]interface{}{
		"status": 0,
		"value":  value,
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func sessionNotCreated(w http.ResponseWriter, err error) {
	body, err := json.Marshal(map[string]interface{}{
		"status":  33,
		"error":   "session not created",
		"message": err.Error(),
		"value": map[string]string{
			"message": err.Error(),
		},
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}

func unknownError(w http.ResponseWriter, err error) {
	body, err := json.Marshal(map[string]interface{}{
		"status":  13,
		"error":   "unknown error",
		"message": err.Error(),
		"value": map[string]string{
			"message": err.Error(),
		},
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}

func unknownMethod(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s is not a supported method for %s", r.Method, r.URL.Path)
	body, err := json.Marshal(map[string]interface{}{
		"status":  9,
		"error":   "unknown method",
		"message": message,
		"value": map[string]string{
			"message": message,
		},
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(body)
}

func unknownCommand(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s %s is an unsupported webdriver command", r.Method, r.URL.Path)
	body, err := json.Marshal(map[string]interface{}{
		"status":  9,
		"error":   "unknown command",
		"message": message,
		"value": map[string]string{
			"message": message,
		},
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusNotFound)
	w.Write(body)
}

func invalidSessionID(w http.ResponseWriter, id string) {
	message := fmt.Sprintf("session %s does not exist", id)
	body, err := json.Marshal(map[string]interface{}{
		"status":  6,
		"error":   "invalid session id",
		"message": message,
		"value": map[string]string{
			"message": message,
		},
	})
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusNotFound)
	w.Write(body)
}
