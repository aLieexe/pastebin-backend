package main

import (
    "log"
    "net/http"
    "os"

    "pastebin-backend/internal/database"
    "pastebin-backend/internal/handlers"
    "pastebin-backend/internal/repository"
)

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}

func main() {
    // connect to database
    if err := database.Connect(); err != nil {
        log.Fatal("failed to connect to database:", err)
    }
    defer database.Close()

    // initialize repository and handlers
    pasteRepo := repository.NewPasteRepository(database.DB)
    pasteHandler := handlers.NewPasteHandler(pasteRepo)

    // setup routes
    http.HandleFunc("/pastes", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            pasteHandler.GetAllPastes(w, r)
        case http.MethodPost:
            pasteHandler.CreatePaste(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    http.HandleFunc("/paste", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            pasteHandler.GetPaste(w, r)
        case http.MethodPut:
            pasteHandler.UpdatePaste(w, r)
        case http.MethodDelete:
            pasteHandler.DeletePaste(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    }))

    // start server
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("failed to start server:", err)
    }
}