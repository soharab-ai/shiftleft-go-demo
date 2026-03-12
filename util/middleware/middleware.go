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

// SanitizeLogInput uses JSON encoding to safely escape all control characters
// This function prevents log forging attacks by encoding user input with proper escaping
// FIXED: Now uses encoding/json for comprehensive control character handling
func SanitizeLogInput(input string) string {
	// Use JSON marshaling to escape all control characters including \n, \r, \t, \x00, 
	// ANSI escape codes, and Unicode line separators (\u2028, \u2029)
	sanitized, err := json.Marshal(input)
	if err != nil {
		// If marshaling fails, return a safe placeholder
		return "[sanitization_error]"
	}
	// Convert JSON string (which includes quotes) to string without quotes
	result := string(sanitized[1 : len(sanitized)-1])
	
	// Limit length to prevent log flooding attacks (additional security measure)
	if len(result) > 200 {
		result = result[:200] + "..."
	}
	return result
}


func LoggingMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		
		// FIXED: Sanitize User-Agent header using comprehensive sanitization
		// and use structured logging with %q format for additional escaping
		userAgent := SanitizeLogInput(r.Header.Get("User-Agent"))
		log.Printf("Request From: user_agent=%q", userAgent)
		
		// FIXED: Sanitize URL path and use structured logging format
		// %q provides Go-syntax quoted string representation with proper escaping
		urlPath := SanitizeLogInput(r.URL.Path)
		log.Printf("Started: method=%s path=%q", r.Method, urlPath)
		
		h(w, r, ps)
		
		// FIXED: Use structured logging with %q format for completion log
		log.Printf("Completed: path=%q duration=%v", urlPath, time.Since(start))
	}
}

            if r != nil {
                switch t := r.(type) {
                case string:
					err = errors.New(t)
                case error:
					err = t
                default:
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
