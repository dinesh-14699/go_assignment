package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"user_service/config"
	"user_service/db"
	grpc_auth "github.com/dinesh-14699/go_assignment/common_utils/grpc_auth"
    cache "github.com/dinesh-14699/go_assignment/common_utils/cache"

	"user_service/internal/handlers"
	"user_service/internal/services"
	"user_service/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)



func main() {
    config.LoadConfig()
    db.InitDB(config.DSN)
    cache.InitializeCache("98.82.179.245:6379", "", 0)

    userService := services.NewUserService(db.DB)
    userHandler := handlers.NewUserHandler(userService)

    var wg sync.WaitGroup

    // Start gRPC server
    wg.Add(1)
    go func() {
        defer wg.Done()
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("Failed to listen on port 50051: %v", err)
        }

        grpcServer := grpc.NewServer()
        grpc_auth.RegisterAuthServiceServer(grpcServer, middleware.GrpcServer())

        log.Println("User Service (Auth) gRPC server is running on port 50051...")
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Failed to serve gRPC server: %v", err)
        }
    }()

    // Start HTTP server
    wg.Add(1)
    go func() {
        defer wg.Done()
        r := chi.NewRouter()
        // r.Use(middleware.Logger)

        r.Post("/register", userHandler.RegisterUser)
        r.Post("/login", userHandler.LoginUser)

        r.Route("/user", func(r chi.Router) {
            r.Use(middleware.AuthMiddleware)
            r.Get("/{userID}", userHandler.GetUser)
            r.Get("/all", userHandler.GetAllUsers)
        })

        log.Println("Starting HTTP server on :8081...")
        log.Fatal(http.ListenAndServe(":8081", r))
    }()

    wg.Wait() 
}
