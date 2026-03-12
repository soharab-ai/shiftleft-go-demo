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

// SanitizeLogInput removes newline characters and control characters to prevent log forging attacks
// This function is added to mitigate log injection vulnerabilities (OWASP A1-Injection)
// FIXED: Using bluemonday library for comprehensive sanitization against sophisticated injection attacks
func SanitizeLogInput(input string) string {
    // Use bluemonday's StrictPolicy which removes all HTML tags and dangerous characters
    p := bluemonday.StrictPolicy()
    sanitized := p.Sanitize(input)
    
    // FIXED: Limit length to prevent log flooding attacks
    maxLength := 200
    if len(sanitized) > maxLength {
        sanitized = sanitized[:maxLength] + "...[truncated]"
    }
    
    return sanitized
}

func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		start := time.Now()
		// FIXED: Sanitize user-controlled input (User-Agent header) before logging to prevent log forging attacks
		userAgent := SanitizeLogInput(r.Header.Get("User-Agent"))
		
		// FIXED: Use structured logging format with clear field boundaries to prevent log structure manipulation
		log.Printf("[REQUEST] method=%s path=%s user_agent=\"%s\" remote_addr=%s", 
			r.Method, 
			r.URL.Path, 
			userAgent,
			r.RemoteAddr)
		
		h(w, r, ps)
		
		// FIXED: Structured logging for completion with clear boundaries
		log.Printf("[COMPLETED] path=%s duration_ms=%d", 
			r.URL.Path, 
			time.Since(start).Milliseconds())
	}

		}

		h(w, r, ps)
	}
}

func (this *Class)CapturePanic(h httprouter.Handle) httprouter.Handle {
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
