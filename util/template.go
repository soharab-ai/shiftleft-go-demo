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
// CORSMiddleware creates a middleware handler for CORS
func CORSMiddleware() func(http.Handler) http.Handler {
    // Get allowed origins from environment or configuration
    allowedOrigins := getAllowedOrigins()
    
    // Create a new CORS handler with specific options
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   allowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
        // Handle OPTIONS preflight requests
        OptionsPassthrough: false,
        // Function to log rejected requests
        Debug: false,
    })
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check if origin is allowed
            origin := r.Header.Get("Origin")
            if !isOriginAllowed(origin, allowedOrigins) && origin != "" {
                logrus.WithFields(logrus.Fields{
                    "origin": origin,
                    "path":   r.URL.Path,
                    "method": r.Method,
                }).Warn("Rejected CORS request from unauthorized origin")
            }
            
            corsHandler.Handler(next).ServeHTTP(w, r)
        })
    }
}

// RenderAsJson renders data as JSON without handling CORS (now handled by middleware)
func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
    // Only set content type, CORS headers are handled by the middleware
    w.Header().Set("Content-Type", "application/json")
    
    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(b)
}

// Helper function to check if an origin is allowed
func isOriginAllowed(origin string, allowedOrigins []string) bool {
    if origin == "" {
        return false
    }
    
    for _, allowed := range allowedOrigins {
        if origin == allowed {
            return true
        }
    }
    return false
}

// Get allowed origins from environment or default to a preset list
func getAllowedOrigins() []string {
    // Read from environment variable if available
    if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
        return []string{origins}
    }
    
    // Default whitelist of trusted domains
    return []string{
        "https://your-trusted-domain.com",
        "https://another-trusted-domain.com",
    }
}

// Example of how to use the middleware in your main application
/*
func main() {
    router := http.NewServeMux()
    
    // Define your routes
    router.HandleFunc("/api/data", handleData)
    
    // Apply the CORS middleware
    handler := CORSMiddleware()(router)
    
    // Start the server
    http.ListenAndServe(":8080", handler)
}
*/

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
