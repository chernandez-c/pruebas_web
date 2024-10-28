// frontend.go
package main

import (
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "path/filepath"
)

type Response struct {
    Message string
    Count   int
}

func main() {
    // Parse templates only once at startup
    indexTemplate := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
    resultTemplate := template.Must(template.ParseFiles(filepath.Join("templates", "result.html")))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Serve the main page template
        indexTemplate.Execute(w, nil)
    })

    http.HandleFunc("/fetch-data", func(w http.ResponseWriter, r *http.Request) {
        // Call the backend API
        resp, err := http.Get("http://localhost:8081/api/data")
        if err != nil {
            log.Println("Error haciendo una petición a http://localhost:8081/api/data")
            http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        // Parse the backend response
        var data Response
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            http.Error(w, "Failed to parse data", http.StatusInternalServerError)
            return
        }

        // Render the response using the result template
        resultTemplate.Execute(w, data)
    })

    log.Println("Frontend server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
