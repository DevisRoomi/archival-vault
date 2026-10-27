package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"archival-vault/backend/database"
	"archival-vault/backend/handlers"
	"archival-vault/backend/middleware"
)

func main() {
	log.Println("[BOOT] Archival Vault System initializing...")

	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("[FATAL] Database initialization failed: %v", err)
	}

	mux := http.NewServeMux()

	// Public Health & System info
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		middleware.JSONSuccess(w, http.StatusOK, map[string]any{
			"status":    "OPERATIONAL",
			"subsystem": "GOLANG-REST-VUE3-ARCHIVE",
			"version":   "1.0.0-wireframe",
		}, "system nominal")
	})

	// Auth Endpoints
	mux.HandleFunc("POST /api/auth/register", handlers.Register)
	mux.HandleFunc("POST /api/auth/login", handlers.Login)
	mux.HandleFunc("GET /api/auth/me", middleware.AuthRequired(handlers.Me))

	// Records Endpoints
	mux.HandleFunc("GET /api/records", handlers.ListRecords)
	mux.HandleFunc("GET /api/records/{id}", handlers.GetRecord)
	mux.HandleFunc("POST /api/records", middleware.AuthRequired(handlers.CreateRecord))
	mux.HandleFunc("PUT /api/records/{id}", middleware.AuthRequired(handlers.UpdateRecord))
	mux.HandleFunc("DELETE /api/records/{id}", middleware.AuthRequired(handlers.DeleteRecord))

	// Telemetry & Stats
	mux.HandleFunc("GET /api/stats", handlers.GetStats)

	// Static Frontend fallback handler
	distPath := filepath.Join("..", "frontend", "dist")
	if _, err := os.Stat(distPath); err != nil {
		distPath = filepath.Join("..", "public")
	}

	fs := http.FileServer(http.Dir(distPath))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If requesting /api, let mux handle 404
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			middleware.JSONError(w, "endpoint not found", http.StatusNotFound)
			return
		}
		// Try static file, fallback to index.html for SPA router
		fPath := filepath.Join(distPath, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(fPath); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	handler := middleware.CORS(mux)

	log.Printf("[SERVER: OK] Archival Vault listening on http://127.0.0.1:%s", port)
	log.Printf("[REGISTRY] Endpoints: /api/auth/* | /api/records/* | /api/stats | /api/health")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("[FATAL] Server crashed: %v", err)
	}
}
