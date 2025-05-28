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
// GetSanitizedUID gets the cookie value and validates it's a valid integer
func GetSanitizedUID(r *http.Request) (string, error) {
    cookie, err := r.Cookie("Uid")
    if err != nil {
        return "", err
    }
    
    // Enhanced validation - whitelist approach for input validation
    // Only allow numeric characters
    for _, ch := range cookie.Value {
        if ch < '0' || ch > '9' {
            return "", errors.New("invalid user ID format: only numeric values allowed")
        }
    }
    
    // Validate that uid is a number and in an acceptable range
    uid, err := strconv.Atoi(cookie.Value)
    if err != nil || uid <= 0 || uid > 1000000 {  // Added range check for defense in depth
        return "", errors.New("invalid user ID format or range")
    }
    
    return strconv.Itoa(uid), nil
}

// Original GetCookie function kept for compatibility but marked as deprecated
// Deprecated: Use GetSanitizedUID instead for secure cookie handling
func GetCookie(r *http.Request, name string) string {
    cookie, _ := r.Cookie(name)
    return cookie.Value
}

func DeleteCookie(w http.ResponseWriter, cookies []string) {
	for _, name := range cookies {
		cookie := &http.Cookie{
			Name:    name,
			Value:   "",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, cookie)
	}
}
