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

// FIX: Enhanced cookie retrieval with error handling to validate cookie existence
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", errors.New("empty cookie value")
	}
	
	// FIX: Add length validation and character whitelisting at cookie layer for defense-in-depth
	if name == "Uid" {
		matched, _ := regexp.MatchString(`^\d{1,10}$`, cookie.Value)
		if !matched {
			return "", errors.New("invalid cookie format")
		}
	}
	
	return cookie.Value, nil
}
