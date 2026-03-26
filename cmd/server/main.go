package main

import (
	"log"
	"net/http"

	"github.com/AzzurroTech/atp/internal/handlers"
	"github.com/AzzurroTech/atp/pkg/storage"
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
