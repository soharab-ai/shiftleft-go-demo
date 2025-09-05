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

// CORSConfig stores CORS configuration parameters
type CORSConfig struct {
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	AllowCredentials    bool
	MaxAge              int
	AllowWildcardSubdomains bool
}

// CORSMiddleware creates a middleware function that handles CORS headers
// properly before request processing
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	// Convert origin slice to map for faster lookups
	allowedOriginsMap := make(map[string]bool)
	wildcardDomains := make([]string, 0)
	
	for _, origin := range config.AllowedOrigins {
		if strings.HasPrefix(origin, "*.") && config.AllowWildcardSubdomains {
			// Store wildcard domains separately for pattern matching
			wildcardDomains = append(wildcardDomains, strings.TrimPrefix(origin, "*."))
		} else {
			allowedOriginsMap[origin] = true
		}
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				allowed := false
				
				// Check exact matches
				if allowedOriginsMap[origin] {
					allowed = true
				}
				
				// Check wildcard subdomains if exact match not found
				if !allowed && config.AllowWildcardSubdomains && len(wildcardDomains) > 0 {
					originURL, err := url.Parse(origin)
					if err == nil {
						host := originURL.Host
						for _, domain := range wildcardDomains {
							if strings.HasSuffix(host, domain) {
								// Ensure it's a subdomain (has a . before the domain)
								if strings.Contains(strings.TrimSuffix(host, domain), ".") {
									allowed = true
									break
								}
							}
						}
					}
				}
				
				// Set CORS headers only for allowed origins
				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin") // Important for proper caching
					
					if config.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}
					
					// Handle preflight OPTIONS request
					if r.Method == http.MethodOptions {
						w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
						if config.MaxAge > 0 {
							w.Header().Set("Access-Control-Max-Age", string(config.MaxAge))
						}
						w.WriteHeader(http.StatusNoContent)
						return
					}
				}
			}
			
			// Add additional security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			
			next.ServeHTTP(w, r)
		})
	}
}

func RenderAsJson(w http.ResponseWriter, r *http.Request, data interfacenull) {
	// Content type is now set here, CORS headers are handled by middleware
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Example usage in main application setup:
// func setupServer() {
//     corsConfig := CORSConfig{
//         AllowedOrigins: []string{
//             "https://trusted-site.com",
//             "https://another-trusted.com",
//             "*.trusted-domain.org"
//         },
//         AllowedMethods: []string{"GET", "POST", "OPTIONS"},
//         AllowedHeaders: []string{"Content-Type", "Authorization"},
//         AllowCredentials: true,
//         MaxAge: 3600,
//         AllowWildcardSubdomains: true,
//     }
//     
//     handler := http.HandlerFunc(yourHandler)
//     http.Handle("/api/", CORSMiddleware(corsConfig)(handler))
// }


func UnSafeRender(w http.ResponseWriter, name string, data ...interface{}) {
	template := template.Must(template.ParseGlob("templates/*"))
	template.ExecuteTemplate(w, name, data)
}

func ToHTML(text string) template.HTML {
	return template.HTML(text)
}
