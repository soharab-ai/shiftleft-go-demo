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
	// MITIGATION: Handle error case when cookie doesn't exist
	if err != nil {
		log.Printf("Cookie retrieval error for '%s': %v", name, err)
		return ""
	}
	
	value := cookie.Value
	// MITIGATION: Remove SQL comment sequences as defense-in-depth measure
	value = strings.ReplaceAll(value, "--", "")
	value = strings.ReplaceAll(value, "/*", "")
	value = strings.ReplaceAll(value, "*/", "")
	
	// Note: Cookie value validation happens at usage point through ValidateUID
	return value
}

