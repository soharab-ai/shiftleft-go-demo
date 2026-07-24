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

func GetCookie(r *http.Request, name string) string {
	cookie, _ := r.Cookie(name)
	return cookie.Value
}
// ValidateUID validates and sanitizes user ID input to prevent SQL injection
// Returns validated integer ID or error if validation fails
func ValidateUID(uid string) (int, error) {
	// FIX: Prevent DoS through excessive input length
	if len(uid) > 10 {
		return 0, errors.New("user ID exceeds maximum length")
	}
	
	// FIX: Use ParseInt with explicit bit size to prevent integer overflow
	id, err := strconv.ParseInt(uid, 10, 32)
	if err != nil || id > 2147483647 {
		return 0, errors.New("invalid user ID format")
	}
	
	// Additional validation: ensure positive ID - FIX: Business logic validation
	if id <= 0 {
		return 0, errors.New("user ID must be positive")
	}
	
	return int(id), nil
}
