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

// GetCookieSafe retrieves and validates cookie value with input sanitization
func GetCookieSafe(r *http.Request, name string) (string, error) {
    cookie, err := r.Cookie(name)
    if err != nil {
        return "", err
    }
    
    // Validate cookie value - for Uid, enforce numeric format to prevent SQL injection
    if name == "Uid" {
        // Enhanced validation: check numeric format and range (Mitigation Note 3)
        uidValue, err := strconv.Atoi(cookie.Value)
        if err != nil {
            return "", fmt.Errorf("invalid cookie format")
        }
        // Whitelist valid user ID range to prevent integer overflow and ensure data integrity
        if uidValue < 1 || uidValue > 2147483647 {
            return "", fmt.Errorf("user ID out of valid range")
        }
    }
    
    return cookie.Value, nil
}
