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
// RateLimiter implements basic rate limiting for preflight requests
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.Mutex
    window   time.Duration
    limit    int
}

// NewRateLimiter creates a new rate limiter with specified window and request limit
func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        window:   window,
        limit:    limit,
    }
}

// Allow checks if the request is allowed based on rate limiting
func (rl *RateLimiter) Allow(ip string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    
    // Clean old requests
    if _, exists := rl.requests[ip]; exists {
        var validRequests []time.Time
        for _, req := range rl.requests[ip] {
            if now.Sub(req) <= rl.window {
                validRequests = append(validRequests, req)
            }
        }
        rl.requests[ip] = validRequests
    }
    
    // Check if rate limit is exceeded
    if len(rl.requests[ip]) >= rl.limit {
        return false
    }
    
    // Add current request
    rl.requests[ip] = append(rl.requests[ip], now)
    return true
}

// Global rate limiter for preflight requests
var preflightLimiter = NewRateLimiter(time.Minute, 60) // 60 requests per minute

// isOriginAllowed checks if the origin is allowed using pattern matching
func isOriginAllowed(origin string, patterns []string) bool {
    if origin == "" {
        return false
    }
    
    for _, pattern := range patterns {
        // Check for subdomain wildcards like "*.example.com"
        if strings.HasPrefix(pattern, "*.") {
            domain := pattern[2:]
            re := regexp.MustCompile(`^https://[^.]+` + regexp.QuoteMeta(domain) + `$`)
            if re.MatchString(origin) {
                return true
            }
        } else if origin == pattern {
            return true
        }
    }
    
    return false
}

// RenderAsJson renders data as JSON with proper CORS security handling
func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
    // Get allowed origins from environment variable, falling back to defaults if not set
    allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
    var allowedOrigins []string
    
    if allowedOriginsEnv != "" {
        // Environment-based configuration for allowed origins
        allowedOrigins = strings.Split(allowedOriginsEnv, ",")
    } else {
        // Default fallback if environment variable is not set
        allowedOrigins = []string{"https://trusted-site.com", "https://another-trusted.com", "*.trusted-domain.com"}
    }
    
    origin := r.Header.Get("Origin")
    
    // Handle preflight OPTIONS requests with rate limiting
    if r.Method == "OPTIONS" {
        clientIP := r.RemoteAddr
        
        // Implement rate limiting for CORS preflight requests
        if !preflightLimiter.Allow(clientIP) {
            w.WriteHeader(http.StatusTooManyRequests)
            return
        }
        
        // Only process if origin is allowed
        if isOriginAllowed(origin, allowedOrigins) {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }
        
        // Add cache control headers for security
        w.Header().Set("Cache-Control", "no-store")
        w.Header().Set("Pragma", "no-cache")
        
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // For regular requests, check if origin is allowed
    if isOriginAllowed(origin, allowedOrigins) {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
    }
    
    // Add cache control headers for security
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("Pragma", "no-cache")
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
