package middleware

import (
    "encoding/json"
    "log"
    "net/http"
)

type errorResponse struct {
    Error errorDetail `json:"error"`
}

type errorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(errorResponse{
        Error: errorDetail{Code: code, Message: message},
    })
}

func Recoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if recovered := recover(); recovered != nil {
                log.Printf("panic while handling %s %s: %v", r.Method, r.URL.Path, recovered)
                WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
            }
        }()

        next.ServeHTTP(w, r)
    })
}
