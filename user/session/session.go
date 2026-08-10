package session

import (
	"fmt"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/config"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type Self struct{}

func New() *Self {
	return &Self{}
}

var store = sessions.NewCookieStore([]byte(config.Cfg.Sessionkey))

func (self *Self) SetSession(w http.ResponseWriter, r *http.Request, data map[string]string) {
	session, err := store.Get(r, "govwa")

	if err != nil {
		log.Println(err.Error())
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false, //set to false for xss :)
	}

	session.Values["govwa_session"] = true

	//create new session to store on server side
	if data != nil {
		for key, value := range data {
			session.Values[key] = value
		}
	}
	err = session.Save(r, w) //safe session and send it to client as cookie

	if err != nil {
		log.Println(err.Error())
	}
}

func (self *Self) GetSession(r *http.Request, key string) string {
	session, err := store.Get(r, "govwa")

	if err != nil {
		log.Println(err.Error())
		return ""
	}
	data := session.Values[key]
	sv := fmt.Sprintf("%v", data)
	return sv
}

func (self *Self) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "govwa")
	if err != nil {
		// MITIGATION FIX: Validate error isn't exposing sensitive data
		errorMsg := err.Error()
		if strings.Contains(strings.ToLower(errorMsg), "password") || 
		   strings.Contains(strings.ToLower(errorMsg), "token") {
			errorMsg = "session retrieval error (details redacted)"
		}
		
		// MITIGATION FIX: Use centralized sanitization from middleware package
		sanitizedError := sanitizeLogInput(errorMsg)
		// MITIGATION FIX: Structured logging with %q for automatic escaping
		log.Printf("[SESSION_ERROR] operation=get error=%q", sanitizedError)
	}

	session.Options = &sessions.Options{
		MaxAge:   -1,
		HttpOnly: true, // MITIGATION FIX: Changed to true for security (prevents XSS attacks)
	}

	session.Values["govwa_session"] = false
	err = session.Save(r, w)

	if err != nil {
		// MITIGATION FIX: Validate error isn't exposing sensitive data
		errorMsg := err.Error()
		if strings.Contains(strings.ToLower(errorMsg), "password") || 
		   strings.Contains(strings.ToLower(errorMsg), "token") {
			errorMsg = "session save error (details redacted)"
		}
		
		// MITIGATION FIX: Use centralized sanitization from middleware package
		sanitizedError := sanitizeLogInput(errorMsg)
		// MITIGATION FIX: Structured logging with %q for automatic escaping
		log.Printf("[SESSION_ERROR] operation=save error=%q", sanitizedError)
	}

	return
}

