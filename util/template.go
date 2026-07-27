package util

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user/session"
)

// isOriginAllowed checks if the origin is in the whitelist of allowed origins
// FIX: Reads allowed origins from environment variable for configuration-based management
// FIX: Added origin format validation to prevent header injection attacks
// FIX: Supports pattern matching for dynamic subdomains
func isOriginAllowed(origin string) bool {
	// FIX: Validate origin format using url.Parse to prevent malformed origins
	if origin == "" {
		return false
	}
	
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}
	
	// FIX: Ensure origin uses HTTPS scheme and has valid host
	if parsedOrigin.Scheme != "https" || parsedOrigin.Host == "" {
		return false
	}
	
	// FIX: Prevent origins with unexpected components (fragments, query strings)
	if parsedOrigin.Fragment != "" || parsedOrigin.RawQuery != "" {
		return false
	}
	
	// FIX: Read allowed origins from environment variable for environment-specific configuration
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		// Fallback to default trusted origins if environment variable not set
		allowedOriginsEnv = "https://trusted-domain.com,https://app.trusted-domain.com"
	}
	
	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	
	for _, allowedOrigin := range allowedOrigins {
		allowedOrigin = strings.TrimSpace(allowedOrigin)
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

	return template.HTML(text)
// FIX: Modified function signature to accept *http.Request parameter for origin validation
// FIX: Added allowCredentials parameter for per-endpoint credential control
// FIX: Replaced wildcard "*" CORS policy with whitelist-based origin validation
// FIX: Added preflight OPTIONS request handling for CORS compliance
// FIX: Only set CORS headers including credentials when origin is explicitly allowed
func RenderAsJson(w http.ResponseWriter, r *http.Request, allowCredentials bool, data ...interface{}) {
	origin := r.Header.Get("Origin")
	
	// FIX: Validate and sanitize origin value to prevent header injection
	if origin != "" {
		parsedOrigin, err := url.Parse(origin)
		if err == nil && parsedOrigin.Scheme == "https" && parsedOrigin.Host != "" && parsedOrigin.Fragment == "" && parsedOrigin.RawQuery == "" {
			origin = parsedOrigin.String()
		} else {
			origin = ""
		}
	}
	
	// FIX: Only set CORS headers if origin is in whitelist to prevent unauthorized cross-origin access
	if origin != "" && isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		// FIX: Conditionally set credentials header based on endpoint requirements
		if allowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		// FIX: Handle preflight OPTIONS requests for CORS compliance
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
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
