// Package-level variables for CORS configuration
var (
	allowedOrigins        map[string]bool
	allowedOriginPatterns []string
	originValidationCache sync.Map
	configMutex           sync.RWMutex
)

// InitializeAllowedOrigins loads allowed origins from environment configuration
// Fixed: Externalized origin allowlist to configuration for operational flexibility
func InitializeAllowedOrigins() {
	configMutex.Lock()
	defer configMutex.Unlock()

	allowedOrigins = make(map[string]bool)
	allowedOriginPatterns = []string{}

	// Fixed: Load exact match origins from environment variable
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins[origin] = true
			}
		}
	}
	// Fixed: Load pattern-based origins for subdomain matching from environment
	patternsEnv := os.Getenv("CORS_ALLOWED_PATTERNS")
	if patternsEnv != "" {
		patterns := strings.Split(patternsEnv, ",")
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				allowedOriginPatterns = append(allowedOriginPatterns, pattern)
			}
		}
	}

	// Default configuration if no environment variables set
	if len(allowedOrigins) == 0 && len(allowedOriginPatterns) == 0 {
		allowedOrigins["https://trusted-domain.com"] = true
		allowedOrigins["https://app.yourdomain.com"] = true
	}
}

}

// validateOrigin checks if an origin is allowed based on exact match or pattern matching
// Fixed: Implements both exact matching and subdomain pattern matching with caching
func validateOrigin(origin string) bool {
	// Fixed: Handle null and missing origins
	if origin == "" || origin == "null" {
		return false
	}

	// Fixed: Validate origin format - must start with https:// for security
	if !strings.HasPrefix(origin, "https://") && !strings.HasPrefix(origin, "http://localhost") {
		return false
	}

	// Fixed: Check cache first for performance optimization
	if cached, found := originValidationCache.Load(origin); found {
		return cached.(bool)
	}

	configMutex.RLock()
	defer configMutex.RUnlock()

	// Fixed: Check exact match in allowlist
	if allowedOrigins[origin] {
		// Cache the validation result with sync.Map for concurrent safety
		originValidationCache.Store(origin, true)
		return true
	}

	// Fixed: Check pattern-based matching for subdomains
	for _, pattern := range allowedOriginPatterns {
		if matchesPattern(origin, pattern) {
			// Cache the validation result
			originValidationCache.Store(origin, true)
			return true
		}
	}

	// Cache negative result to prevent repeated validation attempts
	originValidationCache.Store(origin, false)
	return false
}

// matchesPattern checks if origin matches a wildcard pattern
// Fixed: Pattern matching support for subdomain wildcards
func matchesPattern(origin, pattern string) bool {
	// Support patterns like "https://*.yourdomain.com"
	if strings.HasPrefix(pattern, "https://*.") {
		baseDomain := strings.TrimPrefix(pattern, "https://*.")
// RenderAsJson renders data as JSON with secure CORS headers
// Fixed: Comprehensive CORS security implementation with allowlist validation
func RenderAsJson(w http.ResponseWriter, r *http.Request, data ...interface{}) {
	// Fixed: Validate the Origin header against allowlist with null/empty handling
	origin := r.Header.Get("Origin")
	isOriginAllowed := validateOrigin(origin)

	// Fixed: Handle preflight OPTIONS requests for CORS compliance
	if r.Method == "OPTIONS" {
		if isOriginAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// Fixed: Specify allowed headers for preflight requests
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		// Fixed: Cache preflight response for 24 hours to reduce overhead
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Fixed: Only set CORS headers for explicitly allowed origins (no wildcard)
	if isOriginAllowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// Fixed: Only set credentials header for trusted origins
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Fixed: Maintain existing allowed methods
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	// Fixed: Add additional security headers per OWASP recommendations
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
