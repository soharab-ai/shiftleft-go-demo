package util

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user/session"
)

func SafeRender(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	s := session.New()
	sid := s.GetSession(r, "id") // make uid available to all page
	data["uid"] = sid

	template := template.Must(template.ParseGlob("templates/*"))
	err := template.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err.Error())
	}
}

func RenderAsJson(w http.ResponseWriter, data ...interfacenull) {
	// Define allowed origins instead of using wildcard "*"
	// This should ideally come from configuration
	allowedOrigins := []string{
		"https://yourtrustedomain.com", 
		"https://anothertrustedomain.com",
	}
	
	// Get the origin from the request header
	origin := w.Header().Get("Origin")
	
	// Check if the origin is in our allowed list
	allowOrigin := ""
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			allowOrigin = origin
			break
		}
	}
	
	// Only set CORS headers if origin is allowed
	if allowOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	}
	
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

	w.Write(b)
}

func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
