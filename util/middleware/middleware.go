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

// InitLogger initializes a production-ready structured logger with JSON encoding
// This logger automatically sanitizes all field values to prevent log forging attacks
func InitLogger() *zap.Logger {
    logger, _ := zap.NewProduction()
    return logger
}

var logger = InitLogger()

func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		start := time.Now()
		// Fixed: Using structured logging with zap to prevent log forging attacks (OWASP A1-Injection mitigation)
		// All fields are automatically sanitized and encoded in JSON format, preventing newline and control character injection
		logger.Info("request received",
			zap.String("user_agent", r.Header.Get("User-Agent")),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		
		h(w, r, ps)
		
		// Fixed: Using structured logging for completion message with automatic sanitization
		logger.Info("request completed",
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)))
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
