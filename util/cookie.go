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

// GetCookie retrieves and validates cookie values - FIXED: Added stricter validation for Uid cookie
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	
	// FIXED: Add stricter regex validation for Uid cookie to prevent injection attacks
	if name == "Uid" {
		// FIXED: Regex pattern ensures only positive integers without leading zeros
		matched, _ := regexp.MatchString(`^[1-9][0-9]*$`, cookie.Value)
		if !matched {
			return "", errors.New("invalid cookie value format")
		}
		uid, err := strconv.Atoi(cookie.Value)
		if err != nil || uid <= 0 {
			return "", errors.New("invalid cookie value")
		}
	}
	
	return cookie.Value, nil
}
