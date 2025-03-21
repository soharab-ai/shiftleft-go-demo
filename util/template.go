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
// GetAllowedOrigins returns the list of allowed origins from environment variable
// or falls back to default values if not set
func GetAllowedOrigins() []string {
	// Added dynamic origin configuration from environment variables
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		return strings.Split(origins, ",")
	}
	// Fallback to default trusted origins if not configured
	return []string{"https://trusted-site.com", "https://app.yourcompany.com"}
}

// SetupCORS returns a middleware that handles CORS
func SetupCORS() func(http.Handler) http.Handler {
	// Added dedicated CORS middleware instead of inline implementation
	return cors.New(cors.Options{
		AllowedOrigins:   GetAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours for preflight caching
	}).Handler
}

func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
	origin := r.Header.Get("Origin")
	allowedOrigins := GetAllowedOrigins()
	
	allowOrigin := false
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			allowOrigin = true
			break
		}
	}
	
	// Added handling for preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Only set credentials to true if we've validated the origin
	if allowOrigin {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	} else {
		// Added error logging for rejected origins
		log.Printf("Rejected CORS request from unauthorized origin: %s", origin)
	}
	
	// Added additional security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Usage example:
// func main() {
//     router := http.NewServeMux()
//     router.HandleFunc("/api/data", handleData)
//     
//     // Apply the CORS middleware to all routes
//     http.ListenAndServe(":8080", SetupCORS()(router))
// }
// 
// func handleData(w http.ResponseWriter, r *http.Request) {
//     data := map[string]string{"message": "Hello world"}
//     RenderAsJson(w, r, data)
// }

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
