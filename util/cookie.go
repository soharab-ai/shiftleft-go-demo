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
// Returns validated integer ID or error if input is invalid
func ValidateUID(uid string) (int, error) {
    // Convert to integer and validate - prevents SQL injection by ensuring only numeric values
    id, err := strconv.Atoi(uid)
    if err != nil {
        return 0, errors.New("invalid user ID format")
    }
    
    // Additional validation: ensure positive integer and add upper bound validation to prevent integer overflow and DoS
    if id <= 0 || id > 2147483647 {
        return 0, errors.New("user ID out of valid range")
    }
    
    return id, nil
}
