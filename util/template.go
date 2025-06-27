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
// CORS configuration package
package config

import (
	"os"
	"strings"
)

// GetAllowedOrigins loads allowed origins from environment variables
func GetAllowedOrigins() []string {
	// Default to empty if not configured
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		// Fallback to trusted origins if not configured
		return []string{"https://trusted-site1.com", "https://trusted-site2.com"}
	}
	return strings.Split(originsEnv, ",")
}

// GetCORSMaxAge returns the configured max age for CORS preflight requests
func GetCORSMaxAge() int {
	// Default to 8 hours (in seconds)
	return 28800
}

// Main package with CORS middleware implementation
package util

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	
	"github.com/rs/cors"
	"yourapp/config" // Replace with your actual application's package path
)

// Logger for security events
var securityLogger = log.New(os.Stdout, "SECURITY: ", log.Ldate|log.Ltime|log.Lshortfile)

// ValidateOrigin performs more robust validation of the origin
func ValidateOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	
	// Parse the URL to validate its structure
	u, err := url.Parse(origin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		securityLogger.Printf("Invalid origin format detected: %s", origin)
		return false
	}
	
	// Check against allowed origins
	for _, allowedOrigin := range config.GetAllowedOrigins() {
		if origin == allowedOrigin {
			return true
		}
	}
	
	// Log rejected origins for security monitoring
	securityLogger.Printf("Cross-origin request rejected from: %s", origin)
	return false
}

// CreateCORSHandler creates a CORS middleware with proper configuration
func CreateCORSHandler(specificEndpoint string) *cors.Cors {
	// Allow for endpoint-specific CORS policies
	var allowedOrigins []string
	
	// Example of fine-grained access - different endpoints might have different allowed origins
	if specificEndpoint == "/api/public" {
		allowedOrigins = append(config.GetAllowedOrigins(), "https://public-site.com")
	} else {
		allowedOrigins = config.GetAllowedOrigins()
	}
	
	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           config.GetCORSMaxAge(), // Cache preflight requests to reduce overhead
	})
}

// RenderAsJson safely renders JSON with proper CORS handling
func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
	origin := r.Header.Get("Origin")
	
	// Only set CORS headers if origin is valid and matches allowed origins
	if ValidateOrigin(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Example of middleware usage in your main application setup
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	
	// Add your routes here
	mux.HandleFunc("/api/data", handleData)
	
	// Apply CORS middleware to all routes
	handler := CreateCORSHandler("").Handler(mux)
	
	// For specific endpoints with different CORS policies
	publicHandler := CreateCORSHandler("/api/public").Handler(http.HandlerFunc(handlePublicAPI))
	mux.Handle("/api/public", publicHandler)
	
	return handler
}

// Example handlers
func handleData(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"message": "Data endpoint"}
	RenderAsJson(w, r, data)
}

func handlePublicAPI(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"message": "Public API endpoint"}
	RenderAsJson(w, r, data)
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
