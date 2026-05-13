package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/ashrabya/budget-app/internal/handlers"
	"github.com/ashrabya/budget-app/internal/storage"
)

//go:embed web
var webFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "budget-data.json"
	}

	store, err := storage.NewStore(dataFile)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	h := handlers.New(store)

	mux := http.NewServeMux()

	// API routes
	h.RegisterRoutes(mux)

	// Serve embedded static files
	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("failed to create sub FS: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(webContent)))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Budget App running at http://localhost%s", addr)
	log.Printf("📁 Data stored in: %s", dataFile)

	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
