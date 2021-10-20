package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server is listening...")
	http.HandleFunc("/echo", echoPage)
	http.ListenAndServe(":8080", nil)
}

func echoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Message string `json:"message"`
	}
	var resp struct {
		Echo  string `json:"echo"`
		Error int    `json:"error"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Error = 1
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Echo = req.Message
	json.NewEncoder(w).Encode(resp)
}
