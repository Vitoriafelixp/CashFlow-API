package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
)

const maxJSONBodySize = 1 << 20

func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || (mediaType != "application/json" && !strings.HasSuffix(mediaType, "+json")) {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, maxJSONBodySize+1))
		if err != nil {
			WriteError(w, http.StatusBadRequest, "invalid_json", "Request body could not be read.")
			return
		}
		if len(body) > maxJSONBodySize {
			WriteError(w, http.StatusRequestEntityTooLarge, "payload_too_large", "JSON request body exceeds the 1 MiB limit.")
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if !json.Valid(body) {
			WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must contain valid JSON.")
			return
		}

		next.ServeHTTP(w, r)
	})
}
