package middleware

import (
	"context"
	"strconv"
	"time"
	"user_service/config"
	"github.com/dinesh-14699/common_utils/grpc_auth"

	"github.com/golang-jwt/jwt/v5"
)

type server struct {
	grpc_auth.UnimplementedAuthServiceServer
}

func GrpcServer() *server {
	return &server{}
}

func (s *server) ValidateToken(ctx context.Context, req *grpc_auth.ValidationRequest) (*grpc_auth.ValidationResponse, error) {
	tokenString := req.GetToken()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return &grpc_auth.ValidationResponse{
			IsValid: false,
			Message: "Invalid token",
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &grpc_auth.ValidationResponse{
			IsValid: false,
			Message: "Invalid token claims",
		}, nil
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return &grpc_auth.ValidationResponse{
				IsValid: false,
				Message: "Token has expired",
			}, nil
		}
	} else {
		return &grpc_auth.ValidationResponse{
			IsValid: false,
			Message: "Invalid expiration time in token",
		}, nil
	}

	userID, userOK := claims["user_id"].(float64)
	username, usernameOK := claims["username"].(string)
	email, emailOK := claims["email"].(string)

	if !userOK || !usernameOK || !emailOK {
		return &grpc_auth.ValidationResponse{
			IsValid: false,
			Message: "Invalid user details in token",
		}, nil
	}

	userIDStr := strconv.FormatInt(int64(userID), 10)

	return &grpc_auth.ValidationResponse{
		IsValid:  true,
		Message:  "Token is valid",
		UserId:   userIDStr,
		Username: username,
		Email:     email,
	}, nil
}
