package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"user_service/internal/services"
	"user_service/models"
    
    cache "github.com/dinesh-14699/go_assignment/common_utils/cache"
    "github.com/dinesh-14699/go_assignment/common_utils/logger"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
    Service services.UserServiceInterface
}

func NewUserHandler(service services.UserServiceInterface) *UserHandler {
    return &UserHandler{Service: service}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    logger.Log.Printf("Received request to register user: %v", user.Username)

	if !isValidEmail(user.Email) {
        logger.Log.Printf("Invalid email format for user: %v", user.Username)
        http.Error(w, "Invalid email format", http.StatusBadRequest)
        return
    }

    if err := h.Service.RegisterUser(user); err != nil {
        logger.Log.Printf("Error registering user %v: %v", user.Username, err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    logger.Log.Printf("User registered successfully: %v", user.Username)

    w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        logger.Log.Printf("Error decoding login request: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    logger.Log.Printf("Login attempt for user: %v", user.Username)

    token, err := h.Service.LoginUser(user.Username, user.Password)
    if err != nil {
        logger.Log.Printf("Login failed for user %v: %v", user.Username, err)
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    logger.Log.Printf("User logged in successfully: %v", user.Username)
    w.Write([]byte(token))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userIDStr := chi.URLParam(r, "userID") 
    userID, err := strconv.ParseUint(userIDStr, 10, 32) 
    if err != nil {
        logger.Log.Printf("Invalid user ID: %v", userIDStr)
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    logger.Log.Printf("Fetching user with ID: %v", userID)

    user, err := h.Service.GetUserByID(uint(userID))
    if err != nil {
        logger.Log.Printf("Error fetching user with ID %v: %v", userID, err)
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

	logger.Log.Printf("User fetched successfully with ID: %v", userID)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    logger.Log.Println("Fetching all users...")

	cachedData, err := cache.GetValue("all_users")
	if err == nil && cachedData != "" {
		var cachedUsers []models.User
		err := json.Unmarshal([]byte(cachedData), &cachedUsers)
		if err == nil {
            logger.Log.Println("Cache hit for all users")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedUsers)
			return
		}
        logger.Log.Printf("Error unmarshaling cached data: %v", err)
	}

	users, err := h.Service.GetAllUsers()
	if err != nil {
        logger.Log.Printf("Error fetching all users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    logger.Log.Printf("Fetched %d users from the database", len(users))

	dataToCache, err := json.Marshal(users)
	if err == nil {
		cache.SetValue("all_users", string(dataToCache), 60)
        logger.Log.Println("Users data cached successfully")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}



func isValidEmail(email string) bool {
    regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(email)
}
