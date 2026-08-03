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

// GetCookie retrieves and validates a cookie value from the request
// FIX: Added error handling and validation instead of silently ignoring errors
func GetCookie(r *http.Request, name string) (string, error) {
    cookie, err := r.Cookie(name)
    if err != nil {
        return "", err
    }
    
    // FIX: Validate the cookie value before returning
    if cookie.Value == "" {
        return "", errors.New("empty cookie value")
    }
// ValidateUID validates that the UID is a valid integer and meets security constraints
// FIX: New validation function to ensure UID format is safe before database operations
func ValidateUID(uid string) error {
    // FIX: Ensure uid is a valid integer to prevent SQL injection
    _, err := strconv.Atoi(uid)
    if err != nil {
        return errors.New("invalid user ID format")
    }
    
    // FIX: Additional validation - check reasonable length to prevent overflow attacks
    if len(uid) > 10 {
        return errors.New("user ID exceeds maximum length")
    }
    
    // FIX: Prevent negative or zero values
    if uid[0] == '-' || uid == "0" {
        return errors.New("user ID must be a positive integer")
    }
    
    return nil
}
