package main

import (
	"chatapp/services"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/chat", services.Chat)

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
