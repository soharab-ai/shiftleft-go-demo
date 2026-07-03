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

// Modified to return error when cookie is not found for explicit error handling
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
// ValidateUID validates and sanitizes the user ID input to prevent SQL injection
// Returns validated integer UID and error if validation fails
// Modified to return integer directly and added upper bound validation
func ValidateUID(uid string) (int, error) {
	// Convert to integer and validate - this prevents SQL injection by ensuring only numeric values
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	
	// Additional validation: ensure positive integer and within valid range
	if uidInt <= 0 {
		return 0, errors.New("user ID must be positive")
	}
	
	// Added upper bound check to prevent integer overflow attacks
	if uidInt > 2147483647 {
		return 0, errors.New("user ID out of valid range")
	}
	
	return uidInt, nil
}
