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
// CORSConfig holds all configuration for CORS policy
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedOriginPatterns []*regexp.Regexp
	AllowCredentials bool
	AllowedMethods   []string
	AllowedHeaders   []string
	MaxAge           int
}

// RateLimiter implements a simple rate limiting mechanism
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mu       sync.Mutex
}

var (
	// Global CORS configuration that can be loaded from environment
	corsConfig CORSConfig
	
	// Global rate limiter for cross-origin requests
	rateLimiter *RateLimiter
)

// Initialize CORS configuration from environment or defaults
func init() {
	// Load allowed origins from environment or use defaults
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if originsEnv != "" {
		origins = strings.Split(originsEnv, ",")
	} else {
		origins = []string{"https://trusted-site.com", "https://api.yourservice.com"}
	}
	
	// Compile patterns for wildcard subdomain matching
	patterns := make([]*regexp.Regexp, 0)
	for _, origin := range origins {
		if strings.HasPrefix(origin, "*.") {
			pattern := regexp.MustCompile("^https://[^.]+\\." + regexp.QuoteMeta(strings.TrimPrefix(origin, "*.")) + "$")
			patterns = append(patterns, pattern)
		}
	}
	
	// Parse other CORS settings from environment
	allowCredentials := true
	if creds := os.Getenv("CORS_ALLOW_CREDENTIALS"); creds != "" {
		allowCredentials, _ = strconv.ParseBool(creds)
	}
	
	methods := []string{"GET", "POST", "OPTIONS"}
	if methodsEnv := os.Getenv("CORS_ALLOWED_METHODS"); methodsEnv != "" {
		methods = strings.Split(methodsEnv, ",")
	}
	
	headers := []string{"Content-Type", "Authorization"}
	if headersEnv := os.Getenv("CORS_ALLOWED_HEADERS"); headersEnv != "" {
		headers = strings.Split(headersEnv, ",")
	}
	
	maxAge := 86400 // 24 hours
	if maxAgeEnv := os.Getenv("CORS_MAX_AGE"); maxAgeEnv != "" {
		maxAge, _ = strconv.Atoi(maxAgeEnv)
	}
	
	corsConfig = CORSConfig{
		AllowedOrigins:      origins,
		AllowedOriginPatterns: patterns,
		AllowCredentials:    allowCredentials,
		AllowedMethods:      methods,
		AllowedHeaders:      headers,
		MaxAge:              maxAge,
	}
	
	// Initialize rate limiter (100 requests per minute per origin)
	limit := 100
	if limitEnv := os.Getenv("CORS_RATE_LIMIT"); limitEnv != "" {
		limit, _ = strconv.Atoi(limitEnv)
	}
	
	window := 1 * time.Minute
	rateLimiter = &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// NewCORSMiddleware returns a middleware function that handles CORS
func NewCORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle CORS
			if err := handleCORS(w, r); err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			
			// For preflight OPTIONS requests, we're done
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// Continue with the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// handleCORS applies CORS policy and returns error if request should be rejected
func handleCORS(w http.ResponseWriter, r *http.Request) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Not a CORS request
		return nil
	}
	
	// Check if origin is allowed
	if !isOriginAllowed(origin) {
		return fmt.Errorf("origin not allowed: %s", origin)
	}
	
	// Apply rate limiting
	if !rateLimiter.Allow(origin) {
		return fmt.Errorf("rate limit exceeded for origin: %s", origin)
	}
	
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", origin)
	
	if corsConfig.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))
	}
	
	return nil
}

// isOriginAllowed checks if the origin is in the allowed list or matches a pattern
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	
	// Direct match
	for _, allowed := range corsConfig.AllowedOrigins {
		if allowed == origin || allowed == "*" {
			return true
		}
	}
	
	// Pattern match (e.g., subdomain wildcards)
	for _, pattern := range corsConfig.AllowedOriginPatterns {
		if pattern.MatchString(origin) {
			return true
		}
	}
	
	// Try to parse the URL and check domain
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}
	
	// Extract domain for subdomain matching
	host := parsedOrigin.Host
	for _, allowed := range corsConfig.AllowedOrigins {
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(host, domain) {
				return true
			}
		}
	}
	
	return false
}

// Allow checks if the origin has not exceeded rate limit
func (rl *RateLimiter) Allow(origin string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	
	// Clean up old requests
	if _, exists := rl.requests[origin]; exists {
		var validRequests []time.Time
		for _, t := range rl.requests[origin] {
			if now.Sub(t) <= rl.window {
				validRequests = append(validRequests, t)
			}
		}
		rl.requests[origin] = validRequests
	} else {
		rl.requests[origin] = []time.Timenull
	}
	
	// Check if we're over the limit
	if len(rl.requests[origin]) >= rl.limit {
		return false
	}
	
	// Record this request
	rl.requests[origin] = append(rl.requests[origin], now)
	return true
}

// RenderAsJson renders data as JSON with proper CORS handling
func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
	// CORS handling is now moved to middleware, but we still need to set Content-Type
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
