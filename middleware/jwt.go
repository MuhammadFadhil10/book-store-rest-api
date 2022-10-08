package middleware

import (
	"encoding/json"
	"go-rest/config"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// check jwt auth from cookie
func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type","application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		c, err := r.Cookie("access_token")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				response,_ := json.Marshal(map[string]interface{} {
					"status": false,
					"Error": err.Error(),
				})
				w.Write(response)
				return
			}
		}

		tokenString := c.Value

		claims := &config.JWTClaim{}

		// parse token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			
			err = err.(*jwt.ValidationError)
			response,_ := json.Marshal(map[string]interface{} {
				"status": false,
				"Error": jwt.ValidationErrorSignatureInvalid,
			})
			w.Write(response)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
				response,_ := json.Marshal(map[string]interface{} {
					"status": false,
					"Error": err.Error(),
				})
				w.Write(response)
				return
		}
		

		next.ServeHTTP(w,r)
	})
}