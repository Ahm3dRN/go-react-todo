package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ahm3dRN/go-react-todo/utils"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
type MiddlewareFunc func(http.Handler) http.Handler

// Logging logs all requests with its path and the time it took to process

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// func TokenAuthenticationMiddleware() MiddlewareFunc {

// 	// Create a new Middleware
// 	return func(f http.Handler) http.Handler {

// 		// Define the http.HandlerFunc
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			// Do middleware things
// 			defer func() {
// 				log.Println("that")
// 				token := r.Header.Get("Authentication")
// 				if token == "" {
// 					log.Println(r.URL.Path, r.Method, token)
// 					log.Println("token was not found")
// 				}
// 				user, err := utils.IsValidToken(token)
// 				if err != nil {
// 					log.Println(err)
// 					log.Println("token is invalid")
// 				}
// 				r.Header.Set("UserID", fmt.Sprint(user.ID))

// 			}()

// 			// Call the next middleware/handler in chain
// 			f.ServeHTTP(w, r)
// 		})
// 	}
// }

func TokenAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("that")
		token := r.Header.Get("Authentication")
		user, err := utils.IsValidToken(token)
		if err != nil && r.URL.Path != "/users/login/" && r.URL.Path != "/users/register/" {
			log.Println(r.URL.Path, r.Method, token)
			log.Println("token was not found")
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			en := json.NewEncoder(w)
			en.SetIndent("", "    ")
			en.Encode(map[string]string{"ok": "false", "message": "Athuorization required"})

		} else {
			r.Header.Set("UserID", fmt.Sprint(user.ID))
			next.ServeHTTP(w, r)
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.

	})
}
