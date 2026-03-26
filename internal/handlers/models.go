package models

import "time"

type Source struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type Config struct {
	Version int      `json:"v"`
	Sources []Source `json:"sources"`
}
