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

// Fixed: Enhanced cookie handling to properly return errors
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
// Fixed: Added input validation function with boundary checks to ensure UID is a valid positive integer within reasonable range
func ValidateUID(uid string) (int, error) {
	// Convert to integer to ensure it's numeric
	id, err := strconv.Atoi(uid)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	
	// Fixed: Additional validation - ensure positive integer
	if id <= 0 {
		return 0, errors.New("user ID must be positive")
	}
	
	// Fixed: Add upper boundary check to prevent resource exhaustion and integer overflow
	if id > 2147483647 {
		return 0, errors.New("user ID exceeds maximum allowed value")
	}
	
	return id, nil
}
