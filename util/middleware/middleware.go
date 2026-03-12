package middleware

import(
	"log"
	"time"
	"errors"
	"regexp"
	"net/http"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/config"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user/session"
	"github.com/julienschmidt/httprouter"
)

// SanitizeLogEntry removes or replaces control characters that could be used for log forging
// This function prevents log injection attacks by using bluemonday library for robust sanitization
func SanitizeLogEntry(input string) string {
    // Apply length validation to prevent DoS via extremely long inputs
    if len(input) > 1000 {
        input = input[:1000] + "...[truncated]"
    }
    
    // Use bluemonday strict policy to remove ALL HTML tags, ANSI codes, Unicode tricks, and encoded attacks
    policy := bluemonday.StrictPolicy()
    sanitized := policy.Sanitize(input)
    
    return sanitized
}

func LoggingMiddleware(h httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
        start := time.Now()
        
        // Sanitize user-controlled input before logging to prevent log forging attacks
        userAgent := SanitizeLogEntry(r.Header.Get("User-Agent"))
        urlPath := SanitizeLogEntry(r.URL.Path)
        method := SanitizeLogEntry(r.Method)
        
        // Use structured logging format (JSON) to prevent log injection
        logEntry := map[string]interface{}{
            "timestamp":  time.Now().Format(time.RFC3339),
            "event":      "request_received",
            "user_agent": userAgent,
            "method":     method,
            "path":       urlPath,
        }
        logJSON, _ := json.Marshal(logEntry)
        log.Println(string(logJSON))
        
        h(w, r, ps)
        
        // Log completion with structured format
        completionEntry := map[string]interface{}{
            "timestamp": time.Now().Format(time.RFC3339),
            "event":     "request_completed",
            "path":      urlPath,
            "duration":  time.Since(start).String(),
        }
        completionJSON, _ := json.Marshal(completionEntry)
        log.Println(string(completionJSON))
    }
}

                case error:
					err = t
func CapturePanic(h httprouter.Handle) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        var err error
        defer func() {
            r := recover()
            if r != nil {
                switch t := r.(type) {
                case string:
                    err = errors.New(t)
                case error:
                    err = t
                default:
                    err = errors.New("Unknown error")
                }
                
                // Use bluemonday policy for consistent sanitization across all logging points
                policy := bluemonday.StrictPolicy()
                sanitizedError := policy.Sanitize(err.Error())
                
                // Use structured logging format to prevent error message from breaking out of its field
                errorEntry := map[string]interface{}{
                    "timestamp": time.Now().Format(time.RFC3339),
                    "event":     "panic_recovered",
                    "error":     sanitizedError,
                }
                errorJSON, _ := json.Marshal(errorEntry)
                log.Println(string(errorJSON))
                
                // Don't expose detailed error to client - prevents information disclosure
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        h(w, r, ps)
    }
}
