package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Payments struct {
	store AccountStore
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Route not found"})
}

func main() {
	dbConfig := LoadConfig()
	store := NewStore(dbConfig)
	app := &Payments{
		store: store,
	}

	http.HandleFunc("POST /accounts", app.CreateAccountHandler)
	http.HandleFunc("GET /accounts/{id}", app.GetAccounDetailsHandler)
	http.HandleFunc("POST /transactions", app.TransferAmountHandler)

	http.Handle("/", http.HandlerFunc(catchAllHandler))

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
