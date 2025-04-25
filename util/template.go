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
// Fixed function signature by using interfacenull instead of interfacenull
func RenderAsJson(w http.ResponseWriter, data ...interfacenull) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// CORS handling moved to separate middleware as recommended in mitigation notes
func CorsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	// Convert slice to map for O(1) lookup
	originsMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originsMap[origin] = true
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Only set CORS headers for whitelisted origins - rejecting all others
			if originsMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				
				// Handle preflight OPTIONS requests
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			
			// Proceed with the next handler if origin is allowed or deny implicitly
			if originsMap[origin] {
				next.ServeHTTP(w, r)
			} else {
				// Reject requests from non-whitelisted domains
				http.Error(w, "Origin not allowed", http.StatusForbidden)
			}
		})
	}
}

// Example usage in main application setup
func SetupRoutes() http.Handler {
	// Define trusted origins
	trustedOrigins := []string{
		"https://trusted-site.com",
		"https://another-trusted-site.com",
	}
	
	// Create router/mux
	mux := http.NewServeMux()
	
	// Register your handlers
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Handle API request and use RenderAsJson for response
		data := map[string]string{"message": "Success"}
		RenderAsJson(w, data)
	})
	
	// Wrap with CORS middleware
	return CorsMiddleware(trustedOrigins)(mux)
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
