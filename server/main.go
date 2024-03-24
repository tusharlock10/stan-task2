package main

import (
	"chatapp/services"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	services.InitMonitoring()
	err := services.InitMongoDB()
	if err != nil {
		return
	}

	http.HandleFunc("/messages", services.GetMessagesHandler) // GET endpoint
	http.HandleFunc("/chat", services.ChatHandler)            // WS
	http.Handle("/metrics", promhttp.Handler())               // Monitoring

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
