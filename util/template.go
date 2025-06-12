package util

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user/session"
)

func SafeRender(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	s := session.New()
	sid := s.GetSession(r, "id") // make uid available to all page
	data["uid"] = sid

	template := template.Must(template.ParseGlob("templates/*"))
	err := template.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err.Error())
	}
}

func RenderAsJson(w http.ResponseWriter, data ...interface{}) {
// RenderAsJson serializes data as JSON and writes it to the response.
// CORS handling is now delegated to the middleware.
func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
    // Only set Content-Type header, CORS headers are handled by middleware
    w.Header().Set("Content-Type", "application/json")
    
    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(b)
}

// CORSMiddleware returns a middleware handler that handles CORS requests
// This implements all suggestions from mitigation notes
func CORSMiddleware() func(http.Handler) http.Handler {
    // Get allowed origins from environment or use defaults
    allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
    allowedOrigins := []string{"https://trusted-site.com", "https://app.yourcompany.com"}
    if allowedOriginsEnv != "" {
        // Added environment configuration for origins
        allowedOrigins = strings.Split(allowedOriginsEnv, ",")
    }

    // Using rs/cors library as recommended for robust CORS handling
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   allowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Added explicit allowed headers
        AllowCredentials: true,
        MaxAge:           86400, // 24 hours for preflight requests
        // Added logging for CORS violations
        Debug: true,
        Log: log.Default(),
    })
    
    return corsHandler.Handler
}

// Example of how to use the middleware in your main application:
/*
func SetupServer() http.Handler {
    mux := http.NewServeMux()
    
    // Add your routes
    mux.HandleFunc("/api/data", apiHandler)
    
    // Wrap with CORS middleware
    handler := CORSMiddleware()(mux)
    
    return handler
}
*/

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
