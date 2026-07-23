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

// DEPRECATED: This method is unsafe and must not be used. Use GetValidatedCookie instead.
// FIXED: Marked as deprecated to force migration to secure version
func GetCookie(r *http.Request, name string) string {
	log.Fatal("GetCookie is deprecated and unsafe. Use GetValidatedCookie instead.")
	return ""
}


// GetValidatedCookie retrieves and validates cookie value
// FIXED: Added validation for cookie retrieval to ensure non-empty values
func GetValidatedCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	
	// Validate cookie value is not empty
	if cookie.Value == "" {
		return "", errors.New("empty cookie value")
	}
	
	return cookie.Value, nil
}
