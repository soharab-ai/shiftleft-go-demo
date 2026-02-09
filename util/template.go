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

// Configuration service to manage CORS settings
type CORSConfig struct {
	AllowedOrigins []string
}

// LoadCORSConfig loads allowed origins from environment or defaults
func LoadCORSConfig() *CORSConfig {
	// Load from environment variable or use defaults
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	
	if originsEnv != "" {
		origins = strings.Split(originsEnv, ",")
	} else {
		// Default allowed origins if not specified in environment
		origins = []string{"https://trusted-site.com", "https://admin.trusted-site.com"}
	}
	
	return &CORSConfig{
		AllowedOrigins: origins,
	}
}

// isAllowedOrigin checks if the origin is in the allowed list, supporting wildcards
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
		
		// Support for wildcard subdomains (*.trusted-site.com)
		if strings.HasPrefix(allowed, "*.") {
			domainSuffix := allowed[1:] // remove the *
			if strings.HasSuffix(origin, domainSuffix) {
				return true
			}
		}
	}
	
	return false
}

// CORSMiddleware handles CORS-related security controls
func CORSMiddleware(config *CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			if isAllowedOrigin(origin, config.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			} else if origin != "" {
				// Log rejected origins for security monitoring
				log.Printf("Rejected CORS request from unauthorized origin: %s", origin)
			}
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RenderAsJson now focuses only on JSON rendering without CORS logic
func RenderAsJson(w http.ResponseWriter, data ...interfacenull) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Example of how to use the middleware in your application setup:
/*
func SetupRouter() http.Handler {
    router := mux.NewRouter()
    
    // Load CORS configuration
    corsConfig := LoadCORSConfig()
    
    // Apply CORS middleware
    handler := CORSMiddleware(corsConfig)(router)
    
    // Define routes
    router.HandleFunc("/api/data", HandleData).Methods("GET")
    
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
