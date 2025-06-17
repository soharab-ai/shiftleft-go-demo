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
// Global variables for CORS configuration
var allowedOrigins map[string]bool
var allowedOriginPatterns map[string]*regexp.Regexp

func init() {
	// Initialize allowed origins from environment variables - Implementation of suggestion #1
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	allowedOrigins = make(map[string]bool)
	allowedOriginPatterns = make(map[string]*regexp.Regexp)
	
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins[origin] = true
			
			// Create pattern for subdomain matching - Implementation of suggestion #3
			domain := strings.Replace(strings.Replace(origin, "https://", "", 1), "http://", "", 1)
			pattern := "^https?://([a-zA-Z0-9-]+\\.)*" + regexp.QuoteMeta(domain) + "$"
			allowedOriginPatterns[domain], _ = regexp.Compile(pattern)
		}
	}
	
	// Add default trusted origins if environment variable is not set
	if len(allowedOrigins) == 0 {
		allowedOrigins["https://trusted-site.com"] = true
		allowedOrigins["https://your-app.com"] = true
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing - Implementation of suggestion #2
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks if the origin is in the allowed list or matches allowed patterns
func isOriginAllowed(origin string) bool {
	// Direct match check
	if allowedOrigins[origin] {
		return true
	}
	
	// Pattern match check for subdomains
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	
	for _, pattern := range allowedOriginPatterns {
		if pattern.MatchString(origin) {
			return true
		}
	}
	
	return false
}

// RenderAsJson focuses only on JSON rendering - Implementation of suggestion #4
func RenderAsJson(w http.ResponseWriter, r *http.Request, data ...interfacenull) {
	// Content type setting remains in the render function
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
