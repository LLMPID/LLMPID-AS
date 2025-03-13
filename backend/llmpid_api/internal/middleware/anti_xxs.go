package middleware

import (
	"encoding/json"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
)

func XSSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creates a recorder that acts as a "dummy" ResponseWriter in order to access the body of the response.
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		body, _ := io.ReadAll(rec.Body)

		var jsonData map[string]interface{}
		var jsonDataArray []map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			// If response is not valid JSON, forward it as-is (probably a system error).
			if err := json.Unmarshal(body, &jsonDataArray); err != nil {
				w.WriteHeader(rec.Code)
				w.Write(body)
				return
			}

		}

		// The response body with escaped JS and HTML parameters for XSS remediation.
		var escapedBody []byte

		// Checks if the server response is an array of classification logs or a single one.
		if len(jsonDataArray) != 0 {
			for _, value := range jsonDataArray {
				if _, exists := value["request_text"]; exists {
					value["request_text"] = html.EscapeString(value["request_text"].(string))
				}
			}
			escapedBody, _ = json.Marshal(jsonDataArray)
		} else {

			if _, exists := jsonData["request_text"]; exists {
				jsonData["request_text"] = html.EscapeString(jsonData["request_text"].(string))
			}
			escapedBody, _ = json.Marshal(jsonData)

		}

		// Re-maps headers from the dummy response to the real one.
		for k, v := range rec.Header() {
			for _, sv := range v {
				// Check if header is multi-element and if it was already inserted.
				if r.Header.Get(k) != "" {
					w.Header().Add(k, sv)
				} else {
					w.Header().Set(k, sv)
				}
			}
		}

		// Enforces baseline security headers.
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		w.WriteHeader(rec.Code)
		w.Write(escapedBody)
	})
}
