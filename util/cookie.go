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
	// Return empty string if cookie doesn't exist instead of panicking
	if err != nil {
// IDOR FIX: Rate limiting to prevent enumeration attacks
var (
	rateLimitMutex sync.RWMutex
	rateLimitMap   = make(map[string][]time.Time) // identifier -> list of timestamps
)

func checkRateLimit(identifier string, action string) bool {
	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()
	
	// Create unique key for rate limiting
	key := identifier + ":" + action
	
	// Rate limit configuration: 10 requests per minute
	maxRequests := 10
	timeWindow := 1 * time.Minute
	
	now := time.Now()
	
	// Get existing timestamps for this identifier
	timestamps, exists := rateLimitMap[key]
	if !exists {
		timestamps = []time.Time{}
	}
	
	// Remove timestamps outside the time window
	validTimestamps := []time.Time{}
	for _, ts := range timestamps {
		if now.Sub(ts) < timeWindow {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	
	// Check if rate limit exceeded
	if len(validTimestamps) >= maxRequests {
		return false
	}
	
	// Add current timestamp
	validTimestamps = append(validTimestamps, now)
	rateLimitMap[key] = validTimestamps
	
	return true
}
