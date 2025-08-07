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
func RenderAsJson(w http.ResponseWriter, r *http.Request, data ...interfacenull) {
    // Get origin from request header
    origin := r.Header.Get("Origin")
    
    // Load allowed origins from configuration instead of hardcoding
    // This enables environment-based configuration
    allowedOrigins := config.GetAllowedOrigins()
    
    // Get request method for more granular control
    requestMethod := r.Method
    requestHeaders := r.Header.Get("Access-Control-Request-Headers")
    
    // Check if origin is allowed using pattern matching
    if isOriginAllowed(origin, allowedOrigins) {
        // Set CORS headers with granular method control based on origin
        allowedMethods := cors.GetAllowedMethodsForOrigin(origin)
        allowedHeaders := cors.GetAllowedHeadersForOrigin(origin)
        
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
        
        // Validate and set requested headers if they're allowed
        if requestHeaders != "" && cors.AreHeadersAllowed(requestHeaders, allowedHeaders) {
            w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
        }
        
        // Add Vary header for proper caching
        w.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
    } else if origin != "" {
        // Enhanced error handling with logging for rejected CORS requests
        log.Printf("Rejected CORS request from unauthorized origin: %s", origin)
    }
    
    w.Header().Set("Content-Type", "application/json")
    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(b)
}

// Helper function to check if an origin is allowed using pattern matching
func isOriginAllowed(origin string, allowedPatterns []string) bool {
    if origin == "" {
        return false
    }
    
    // First check for exact matches
    for _, pattern := range allowedPatterns {
        if pattern == origin {
            return true
        }
        
        // Support wildcard subdomain matching (e.g., *.example.com)
        if strings.HasPrefix(pattern, "*.") {
            suffix := pattern[1:] // remove the *
            if strings.HasSuffix(origin, suffix) {
                // Ensure it's a proper subdomain
                domainPart := origin[:len(origin)-len(suffix)]
                if !strings.Contains(domainPart, ".") {
                    return true
                }
            }
        }
    }
    
    // Check for regex patterns if exact matching fails
    for _, pattern := range allowedPatterns {
        if strings.HasPrefix(pattern, "regex:") {
            regexPattern := pattern[6:] // Remove "regex:" prefix
            matched, err := regexp.MatchString(regexPattern, origin)
            if err == nil && matched {
                return true
            }
        }
    }
    
    // Check dynamic whitelist from database or service if configured
    if cors.IsDynamicWhitelistEnabled() {
        return cors.CheckDynamicWhitelist(origin)
    }
    
    return false
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
