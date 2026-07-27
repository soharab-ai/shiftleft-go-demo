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

func RenderAsJson(w http.ResponseWriter, r *http.Request, data ...interface{}) {
	// FIX: Load allowed origins from environment variable for configuration-based management
	allowedOriginsEnv := os.Getenv("ALLOWED_CORS_ORIGINS")
	var allowedOriginsSlice []string
	
	if allowedOriginsEnv != "" {
		allowedOriginsSlice = strings.Split(allowedOriginsEnv, ",")
		// Trim whitespace from each origin
		for i, origin := range allowedOriginsSlice {
			allowedOriginsSlice[i] = strings.TrimSpace(origin)
		}
	}
	
	// FIX: Convert slice to map for O(1) lookup performance
	allowedOriginsMap := make(map[string]bool)
	for _, origin := range allowedOriginsSlice {
		// FIX: Validate origin format and ensure only HTTPS origins are allowed
		if parsedURL, err := url.Parse(origin); err == nil && parsedURL.Scheme == "https" {
			// FIX: Store in lowercase for case-insensitive comparison
			allowedOriginsMap[strings.ToLower(origin)] = true
		}
	}
	
	// FIX: Extract and validate the request origin against the whitelist
	origin := r.Header.Get("Origin")
	isAllowed := false
	
	if origin != "" {
		// FIX: Validate origin format before comparison
		if parsedOrigin, err := url.Parse(origin); err == nil && parsedOrigin.Scheme == "https" {
			// FIX: Case-insensitive comparison to prevent bypass
			isAllowed = allowedOriginsMap[strings.ToLower(origin)]
		}
	}
	
	// FIX: Handle OPTIONS preflight requests for complex CORS requests
	if r.Method == http.MethodOptions {
		if isAllowed {
			// FIX: Set specific origin instead of wildcard for preflight
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// FIX: Set allowed methods only for whitelisted origins
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			// FIX: Cache preflight response for 24 hours to reduce request frequency
			w.Header().Set("Access-Control-Max-Age", "86400")
			// FIX: Inform caching proxies that responses vary based on Origin header
			w.Header().Set("Vary", "Origin")
		}
		// Return 204 No Content for preflight requests
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	// FIX: Only set CORS headers if origin is whitelisted to prevent CSRF attacks
	if isAllowed {
		// FIX: Set specific origin instead of wildcard
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// FIX: Set allowed methods only for whitelisted origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		// FIX: Prevent cache poisoning attacks with Vary header
		w.Header().Set("Vary", "Origin")
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
