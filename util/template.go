// Package-level cache for allowed origins loaded from configuration
var allowedOriginsCache []string

// init loads and validates allowed origins from environment configuration
// FIXED: Moved hardcoded origins to environment-based configuration
func init() {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv == "" {
		// Default to empty list - no origins allowed if not configured
		allowedOriginsCache = []string{}
		return
	}
	
	// Parse comma-separated list and validate each origin
// isOriginAllowed validates if the requesting origin is in the allowlist
// FIXED: Added origin validation with subdomain pattern matching support
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	
	// FIXED: Check against configuration-based cache instead of hardcoded list
	for _, allowed := range allowedOriginsCache {
		// FIXED: Support wildcard subdomain patterns (e.g., *.yourdomain.com)
		if strings.HasPrefix(allowed, "*.") {
			// Extract the domain pattern (everything after "*")
			domainPattern := allowed[1:] // Results in ".yourdomain.com"
			
			// Ensure proper subdomain matching with dot separator
			// This prevents matching "malicious-yourdomain.com"
			if strings.HasSuffix(origin, domainPattern) {
				// Additional validation: ensure it's a proper subdomain match
				if matched, _ := regexp.MatchString(`^https?://[a-zA-Z0-9\-]+`+regexp.QuoteMeta(domainPattern)+"$", origin); matched {
					return true
				}
			}
		} else if origin == allowed {
			// FIXED: Exact match for non-wildcard origins
			return true
		}
	}
	return false
}

	w.Write(b)
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
// FIXED: Updated RenderAsJson with method-based controls, preflight handling, and security logging
// This prevents unauthorized cross-origin access and CSRF attacks
func RenderAsJson(w http.ResponseWriter, r *http.Request, allowedMethods string, data ...interface{}) {
	origin := r.Header.Get("Origin")
	
	// FIXED: Only set CORS headers if origin is in the allowlist (replaces wildcard "*")
	if isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// FIXED: Use method-based CORS controls instead of hardcoded "POST, GET"
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		
		// FIXED: Handle OPTIONS preflight requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	} else {
		// FIXED: Add security logging for rejected origins to enable monitoring
		if origin != "" {
			// Use structured logging with separate fields to prevent log injection
			log.Printf("Rejected cross-origin request | origin=%s | path=%s | remote_addr=%s", 
				sanitizeLogField(origin), 
				sanitizeLogField(r.URL.Path), 
// sanitizeLogField prevents log injection by removing newlines and control characters
// FIXED: Added helper function for safe logging to prevent log forging attacks
func sanitizeLogField(field string) string {
	// Remove newlines, carriage returns, and other control characters
	sanitized := strings.ReplaceAll(field, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	sanitized = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(sanitized, "")
	return sanitized
}

		return
	}
	w.Write(b)
}
