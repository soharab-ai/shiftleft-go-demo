package util

import (
	"net/http"
	"time"
)

func SetCookieLevel(w http.ResponseWriter, r *http.Request, cookievalue string) {

	level := cookievalue
	if level == "" {
		level = "low"
	}
	SetCookie(w, "Level", level)

}

func CheckLevel(r *http.Request) bool {
	level := GetCookie(r, "Level")
	if level == "" || level == "low" {
		return false //set default level to low
	} else if level == "high" {
		return true //level == high
	} else {
		return false // level == low
	}
}

/* cookie setter getter */

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		//Path : "/",
		//Domain : "localhost",
		Name:  name,
		Value: value,
	}
	http.SetCookie(w, &cookie)
}

// GetCookie retrieves and validates cookie value with enhanced input sanitization
func GetCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	
	// Validate format - only digits allowed (mitigation: format validation)
	matched, _ := regexp.MatchString(`^[0-9]+$`, cookie.Value)
	if !matched {
		log.Printf("Invalid cookie format for %s from IP: %s", name, sanitizeIPForLog(r.RemoteAddr))
		return ""
	}
	
	// Additional validation: length check to prevent overflow (mitigation: length restriction)
	if len(cookie.Value) > 10 {
		log.Printf("Cookie value exceeds maximum length for %s", name)
		return ""
	}
	
	// Range validation to ensure valid database integer bounds (mitigation: range validation)
	uid, err := strconv.ParseInt(cookie.Value, 10, 32)
	if err != nil || uid < 1 || uid > 2147483647 {
		log.Printf("Cookie value out of valid range for %s", name)
		return ""
	}
	
	return cookie.Value
}

// sanitizeIPForLog removes potentially dangerous characters from IP addresses for safe logging
func sanitizeIPForLog(ip string) string {
	// Remove port and sanitize for log injection prevention (mitigation: log injection prevention)
	return regexp.MustCompile(`[^0-9.:a-fA-F]`).ReplaceAllString(ip, "")
}
