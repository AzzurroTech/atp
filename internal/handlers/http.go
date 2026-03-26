package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AzzurroTech/atp/internal/models"
	"github.com/AzzurroTech/atp/pkg/storage"
)

type Handler struct {
	store storage.Store
}

func RegisterRoutes(store storage.Store) {
	h := &Handler{store: store}
	http.HandleFunc("/", h.handleRoot)
	http.HandleFunc("/add", h.handleAdd)
	http.HandleFunc("/api/export", h.handleExport)
	http.HandleFunc("/api/import", h.handleImport)
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	sources := h.store.GetAll()
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].CreatedAt.After(sources[j].CreatedAt)
	})

	tmpl, err := template.New("index").Funcs(template.FuncMap{
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
		},
	}).Parse(IndexHTML)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, map[string]interface{}{"Sources": sources})
	if err != nil {
		http.Error(w, "Template execution failed", 500)
	}
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	urlStr := r.FormValue("url")
	notes := r.FormValue("notes")

	if urlStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	if h.store.Exists(urlStr) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newSource := &models.Source{
		ID:        time.Now().Format("20060102150405"),
		URL:       urlStr,
		Title:     "New Source: " + urlStr,
		Summary:   "Summary generated for " + urlStr + ". (Simulated)",
		Notes:     notes,
		CreatedAt: time.Now(),
	}

	h.store.Add(newSource)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) handleExport(w http.ResponseWriter, r *http.Request) {
	sources := h.store.GetAll()
	cfg := models.Config{Version: 1, Sources: sources}

	jsonData, _ := json.Marshal(cfg)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonData)
	gz.Close()

	encoded := base64.URLEncoding.EncodeToString(buf.Bytes())
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(encoded))
}

func (h *Handler) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct{ Config string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	data, err := base64.URLEncoding.DecodeString(req.Config)
	if err != nil {
		http.Error(w, "Invalid Base64", 400)
		return
	}

	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Invalid GZIP data", 400)
		return
	}
	defer gr.Close()

	decompressed, _ := io.ReadAll(gr)
	var cfg models.Config
	if err := json.Unmarshal(decompressed, &cfg); err != nil {
		http.Error(w, "Invalid JSON structure", 400)
		return
	}

	for _, s := range cfg.Sources {
		if !h.store.Exists(s.URL) {
			h.store.Add(&s)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
