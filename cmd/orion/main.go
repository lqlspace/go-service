package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	dataStore = make(map[int]Item)
	counter   = 1
)

func main() {
	http.HandleFunc("/api/items", handleItems)
	http.HandleFunc("/api/items/", handleItemByID)
	http.HandleFunc("/api/health", healthCheck)

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	log.Printf("Starting server on %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var items []Item
		for _, v := range dataStore {
			items = append(items, v)
		}
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		item.ID = counter
		counter++
		dataStore[item.ID] = item
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleItemByID(w http.ResponseWriter, r *http.Request) {
	id := 0
	if _, err := fmt.Sscanf(r.URL.Path, "/api/items/%d", &id); err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, exists := dataStore[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(item)
	case http.MethodDelete:
		delete(dataStore, id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}