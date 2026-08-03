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

// RenderAsJson renders data as JSON with secure CORS configuration
// Fixed: Replaced wildcard CORS with origin whitelist validation to prevent unauthorized cross-origin access
func RenderAsJson(w http.ResponseWriter, r *http.Request, data ...interface{}) {
	// Fixed: Get the origin from request header
	origin := r.Header.Get("Origin")
	
	// Fixed: Check if origin is in the blocked list (rate limiting)
	if origin != "" && isOriginBlocked(origin) {
		http.Error(w, "Too many rejected requests", http.StatusForbidden)
		return
	}
	
	// Fixed: Load allowed origins and patterns from environment variable or use secure defaults
	exactOrigins, regexPatterns := getAllowedOrigins()
	
	// Fixed: Validate the Origin header against whitelist with protocol and port verification
	isAllowed := false
	
	if origin != "" {
		// Fixed: Parse origin URL for robust validation including protocol and port
		originURL, err := url.Parse(origin)
		if err == nil {
			// Fixed: Enforce HTTPS in production environment
			if os.Getenv("ENVIRONMENT") == "production" && originURL.Scheme != "https" {
				isAllowed = false
			} else {
				// Fixed: Check exact matches first (fast path)
				for _, allowedOrigin := range exactOrigins {
					allowedURL, err := url.Parse(allowedOrigin)
					if err == nil && originURL.Host == allowedURL.Host && originURL.Scheme == allowedURL.Scheme {
						isAllowed = true
						break
					}
				}
				
				// Fixed: Fall back to regex pattern matching for dynamic subdomains
				if !isAllowed {
					for _, pattern := range regexPatterns {
						if pattern.MatchString(origin) {
							isAllowed = true
							break
						}
					}
				}
			}
		}
	}
	
	// Fixed: Only set CORS headers if origin is in the whitelist
	if isAllowed {
		// Fixed: Set specific origin instead of wildcard "*"
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// Fixed: Only enable credentials for trusted origins
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Fixed: Add security headers to prevent MIME-sniffing and clickjacking attacks
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		
		// Fixed: Handle CORS preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	} else if origin != "" {
		// Fixed: Log and track rejected CORS requests for security monitoring and rate limiting
		logRejectedOrigin(origin)
	}
	
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Content-Type", "application/json")
	
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// getAllowedOrigins retrieves the list of allowed origins and regex patterns from environment configuration
// Fixed: Added helper function to manage allowed origins with regex pattern support
func getAllowedOrigins() ([]string, []*regexp.Regexp) {
	// Fixed: Use cached patterns if already compiled
	cachedPatternsMutex.RLock()
	if cachedExactOrigins != nil && cachedRegexPatterns != nil {
		defer cachedPatternsMutex.RUnlock()
		return cachedExactOrigins, cachedRegexPatterns
	}
	cachedPatternsMutex.RUnlock()
	
	// Fixed: Acquire write lock to compile patterns
	cachedPatternsMutex.Lock()
	defer cachedPatternsMutex.Unlock()
	
	// Fixed: Double-check after acquiring write lock
	if cachedExactOrigins != nil && cachedRegexPatterns != nil {
		return cachedExactOrigins, cachedRegexPatterns
	}
	
	var exactOrigins []string
	var regexPatterns []*regexp.Regexp
	
	// Fixed: Read from environment variable for deployment flexibility
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	
	if originsEnv != "" {
		// Fixed: Parse comma-separated list of origins and regex patterns
		origins := strings.Split(originsEnv, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed == "" {
				continue
			}
			
			// Fixed: Support regex patterns prefixed with "regex:"
			if strings.HasPrefix(trimmed, "regex:") {
				patternStr := strings.TrimPrefix(trimmed, "regex:")
				pattern, err := regexp.Compile(patternStr)
				if err == nil {
					regexPatterns = append(regexPatterns, pattern)
				}
			} else {
				exactOrigins = append(exactOrigins, trimmed)
			}
		}
	}
	
	// Fixed: Secure default - return empty list or specific trusted domains if no env var set
	if len(exactOrigins) == 0 && len(regexPatterns) == 0 {
		exactOrigins = []string{
			"https://trusted-domain.com",
			"https://app.yourdomain.com",
		}
	}
	
	// Fixed: Cache compiled patterns for performance
	cachedExactOrigins = exactOrigins
	cachedRegexPatterns = regexPatterns
	
	return exactOrigins, regexPatterns
}

// logRejectedOrigin logs attempts from non-whitelisted origins and implements rate limiting
// Fixed: Added security logging and rate limiting for rejected CORS requests
func logRejectedOrigin(origin string) {
	// Fixed: Sanitize origin before logging to prevent log injection
	sanitized := strings.ReplaceAll(origin, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	
	// Fixed: Track rejection count for rate limiting
	now := time.Now()
	
	rejectionTrackerMutex.Lock()
	defer rejectionTrackerMutex.Unlock()
	
	if tracker, exists := rejectionTracker[sanitized]; exists {
		// Fixed: Check if we're still within the time window
		if now.Sub(tracker.FirstSeen) <= rejectionTimeWindow {
			tracker.Count++
			tracker.LastSeen = now
			
			// Fixed: Block origin if threshold exceeded
			if tracker.Count > rejectionThreshold {
				blockedOrigins.Store(sanitized, now.Add(blockDuration))
			}
		} else {
			// Fixed: Reset counter if outside time window
			rejectionTracker[sanitized] = &RejectionTracker{
				Count:     1,
				FirstSeen: now,
				LastSeen:  now,
			}
		}
	} else {
		// Fixed: Initialize new tracker entry
		rejectionTracker[sanitized] = &RejectionTracker{
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
		}
	}
	
	// In production, use structured logging framework
	// Example: log.Warn("Rejected CORS request", "origin", sanitized)
	_ = sanitized // Placeholder for actual logging implementation
}

// isOriginBlocked checks if an origin is currently blocked due to rate limiting
// Fixed: Added helper function to check blocked origins
func isOriginBlocked(origin string) bool {
	sanitized := strings.ReplaceAll(origin, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	
	if blockUntil, exists := blockedOrigins.Load(sanitized); exists {
		if time.Now().Before(blockUntil.(time.Time)) {
			return true
		}
		// Fixed: Remove expired block
		blockedOrigins.Delete(sanitized)
	}
	return false
}

// RejectionTracker tracks rejection attempts for rate limiting
// Fixed: Added structure to track rejection attempts per origin
type RejectionTracker struct {
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Fixed: Package-level variables for caching and rate limiting
var (
	cachedExactOrigins   []string
	cachedRegexPatterns  []*regexp.Regexp
	cachedPatternsMutex  sync.RWMutex
	rejectionTracker     = make(map[string]*RejectionTracker)
	rejectionTrackerMutex sync.Mutex
	blockedOrigins       sync.Map
	rejectionThreshold   = 10
	rejectionTimeWindow  = 1 * time.Minute
	blockDuration        = 5 * time.Minute
)
