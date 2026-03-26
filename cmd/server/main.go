package main

import (
	"log"
	"net/http"

	"github.com/azzurrotech/ATP/internal/handlers"
	"github.com/azzurrotech/ATP/pkg/storage"
)

func main() {
	// Initialize in-memory storage
	store := storage.NewMemoryStore()

	// Register handlers
	handlers.RegisterRoutes(store)

	addr := ":8080"
	log.Printf("🚀 Azzurro Technology Platform starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
