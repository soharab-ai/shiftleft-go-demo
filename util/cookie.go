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
    cookie, err := r.Cookie(name)
    if err != nil {
        return ""
    }
    
    // Validate that the value is a positive integer
    id, err := strconv.Atoi(cookie.Value)
    if err != nil || id <= 0 {
        return ""
    }
    
    // Add maximum length check to prevent numeric overflow attacks
    const MAX_SAFE_ID = 1000000000 // Set appropriate limit for your application
    if len(cookie.Value) > 10 || id > MAX_SAFE_ID {
        return ""
    }
    
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
