package middleware

import(
	"log"
	"time"
	"errors"
	"regexp"
	"net/http"

// MITIGATION FIX: Comprehensive control character regex to prevent log forging attacks
var sessionControlCharsRegex = regexp.MustCompile(`[\x00-\x1F\x7F\u2028\u2029]`)
var sessionAnsiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// sanitizeLogInput removes control characters, ANSI escape sequences, and Unicode newlines
// to prevent log forging attacks - MITIGATION: Enhanced with regex-based comprehensive sanitization
func sanitizeLogInput(input string) string {
    // MITIGATION FIX: Remove ALL control characters including Unicode variants
    sanitized := sessionControlCharsRegex.ReplaceAllString(input, "")
    
    // MITIGATION FIX: Remove ANSI escape sequences that can manipulate terminal output
    sanitized = sessionAnsiRegex.ReplaceAllString(sanitized, "")
    
    // MITIGATION FIX: Normalize whitespace
    sanitized = strings.TrimSpace(sanitized)
    
    // MITIGATION FIX: Length limiting with proper UTF-8 handling
    if len([]rune(sanitized)) > 200 {
        runes := []rune(sanitized)
        sanitized = string(runes[:200]) + "...[truncated]"
    }
    return sanitized
}

    return sanitized
}
func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			
			// MITIGATION FIX: Sanitize and structure the error log with request context
			sanitizedError := SanitizeLogInput(err.Error())
			sanitizedPath := SanitizeLogInput(r.URL.Path)
			
			// MITIGATION FIX: Use structured logging with %q verb for automatic escaping
			log.Printf("[PANIC] path=%q error=%q method=%s", 
				sanitizedPath, 
				sanitizedError, 
				r.Method)
			
			// MITIGATION FIX: Generic error response prevents information disclosure
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	h(w, r, ps)
}

	// MITIGATION FIX: Use structured logging with fields instead of string interpolation
	logFields := map[string]interface{}{
		"user_agent": SanitizeLogInput(r.Header.Get("User-Agent")),
		"method":     r.Method,
		"path":       SanitizeLogInput(r.URL.Path),
		"remote_ip":  SanitizeLogInput(r.RemoteAddr),
	}
	
	// MITIGATION FIX: Using %q format verb for automatic escaping as defense-in-depth
	log.Printf("[REQUEST] user_agent=%q method=%s path=%q remote_ip=%q", 
		logFields["user_agent"], 
		logFields["method"], 
		logFields["path"], 
		logFields["remote_ip"])
	
	h(w, r, ps)
	
	// MITIGATION FIX: Structured logging prevents injection even if sanitization has gaps
	log.Printf("[COMPLETE] method=%s path=%q duration=%v", 
		logFields["method"], 
		logFields["path"], 
		time.Since(start))
}

                    err = errors.New("Unknown error")
				}
				log.Println(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r, ps)
	}
}

func(this *Class)DetectSQLMap(h httprouter.Handle)httprouter.Handle{
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userAgent := r.Header.Get("User-Agent")
		sqlmapDetected, _ := regexp.MatchString("sqlmap*", userAgent)
		if sqlmapDetected{
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			log.Printf("sqlmap detect ")
			return
		}else{
			h(w, r, ps)
		}
	}
}
