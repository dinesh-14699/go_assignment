package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
    "github.com/dinesh-14699/go_assignment/common_utils/grpc_auth"
	"google.golang.org/grpc"
)

const UserContextKey = "user_id"

func TokenValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            http.Error(w, "Bearer token missing", http.StatusUnauthorized)
            return
        }

        isValid, err := validateTokenWithGRPC(token)
        fmt.Println(err, isValid, token)
        if err != nil || !isValid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func validateTokenWithGRPC(token string) (bool, error) {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        return false, err
    }
    defer conn.Close()

    client := grpc_auth.NewAuthServiceClient(conn)

    response, err := client.ValidateToken(context.Background(), &grpc_auth.ValidationRequest{Token: token})

    fmt.Print(response, err)

    if err != nil {
        return false, err
    }

    return response.IsValid, nil
}
