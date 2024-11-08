package middleware

import (
    "context"
    "net/http"
    "strings"
    "time"
    "user_service/config"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey = contextKey("user_id")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Bearer token missing", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(config.JWTSecret), nil
        })

        
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if exp, ok := claims["exp"].(float64); ok {
                if time.Now().Unix() > int64(exp) {
                    http.Error(w, "Token has expired", http.StatusUnauthorized)
                    return
                }
            } else {
                http.Error(w, "Invalid expiration time in token", http.StatusUnauthorized)
                return
            }

            if userID, ok := claims["user_id"].(float64); ok {
                ctx := context.WithValue(r.Context(), UserContextKey, uint(userID))
                r = r.WithContext(ctx)
            } else {
                http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
                return
            }
        } else {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
