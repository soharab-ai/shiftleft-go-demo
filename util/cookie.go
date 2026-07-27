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
// Fixed: Added error handling and validation to prevent empty cookie values
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	// Validate cookie value is not empty
	if cookie.Value == "" {
		return "", errors.New("empty cookie value")
// ValidateUID validates and sanitizes the user ID input
// Fixed: Added comprehensive input validation including length check, format validation, and upper bound enforcement
func ValidateUID(uid string) (int, error) {
	// Fixed: Add length check before conversion to prevent DoS from extremely long numeric strings
	if len(uid) > 10 {
		return 0, errors.New("user ID format invalid")
	}
	
	// Convert to integer and validate format
	id, err := strconv.Atoi(uid)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	
	// Additional validation: ensure positive integer
	if id <= 0 {
		return 0, errors.New("user ID must be positive")
	}
	
	// Fixed: Add maximum value validation to prevent integer overflow and enforce reasonable upper bound
	if id > 2147483647 {
		return 0, errors.New("user ID exceeds maximum allowed value")
	}
	
	return id, nil
}
