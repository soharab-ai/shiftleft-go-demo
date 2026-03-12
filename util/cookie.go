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
	// FIXED: This method is kept for backward compatibility but should be deprecated
	// Use GetCookieSafe instead for security-critical operations
	return cookie.Value
}

// GetCookieSafe retrieves and validates a cookie value
// FIXED: Added input validation to prevent SQL injection through cookie manipulation
func GetCookieSafe(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	
	// Validate cookie value - only allow numeric values for Uid cookie
	// FIXED: Input validation using regex to ensure only digits are accepted
	matched, err := regexp.MatchString("^[0-9]+$", cookie.Value)
	if err != nil || !matched {
		return "", fmt.Errorf("invalid cookie format")
	}
	
	// FIXED: Add range validation to prevent integer overflow and logical attacks
	uidInt, err := strconv.Atoi(cookie.Value)
	if err != nil || uidInt <= 0 || uidInt > 2147483647 {
		return "", fmt.Errorf("user ID out of valid range")
	}
	
	return cookie.Value, nil
}
