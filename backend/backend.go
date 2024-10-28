// backend.go
package main
import (
    "io"
    "log"
    "net/http"
    "time"
    "encoding/json"
    "bytes" 
)


type Response struct {
    Message string `json:"message"`
    Count   int    `json:"count"`
}

// LoggingMiddleware logs details of incoming HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        // Set CORS headers to allow requests from http://localhost:8080
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // Handle preflight OPTIONS request for CORS
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Log request method and URL
        log.Printf("Received %s request for %s", r.Method, r.URL)


        // Log request headers
        /*
        for name, values := range r.Header {
            for _, value := range values {
                log.Printf("Header: %s=%s", name, value)
            }
        }
        */

        // Log request body if it exists
        if r.Body != nil {
            body, err := io.ReadAll(r.Body)
            if err == nil {
                log.Printf("Body: %s", string(body))
                
                // Reset the body so other handlers can read it
                r.Body = io.NopCloser(io.Reader(bytes.NewReader(body)))
            } else {
                log.Printf("Error reading body: %v", err)
            }
        }

        // Call the next handler in the chain
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("Completed in %v", duration)
    })
}

func main() {
    mux := http.NewServeMux()

    // Example handler
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })

    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        response := Response{
            Message: "Hello from the backend!",
            Count:   42,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })


    // Wrap the mux with the logging middleware
    loggedMux := LoggingMiddleware(mux)

    // Start the server
    log.Println("Backend Server running on http://localhost:8081")
    log.Fatal(http.ListenAndServe(":8081", loggedMux))
}
