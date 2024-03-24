package services

import (
	"net/http"

	"encoding/json"
)

// Get all the messages stored in the db
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
	}
}
