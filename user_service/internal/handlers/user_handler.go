package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"user_service/internal/services"
	"user_service/models"
    
    cache "github.com/dinesh-14699/go_assignment/common_utils/cache"

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

	if !isValidEmail(user.Email) {
        http.Error(w, "Invalid email format", http.StatusBadRequest)
        return
    }

    if err := h.Service.RegisterUser(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    token, err := h.Service.LoginUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    w.Write([]byte(token))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userIDStr := chi.URLParam(r, "userID") 
    userID, err := strconv.ParseUint(userIDStr, 10, 32) 
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    user, err := h.Service.GetUserByID(uint(userID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	cachedData, err := cache.GetValue("all_users")
	if err == nil && cachedData != "" {
		var cachedUsers []models.User
		err := json.Unmarshal([]byte(cachedData), &cachedUsers)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedUsers)
			return
		}
	}

	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataToCache, err := json.Marshal(users)
	if err == nil {
		cache.SetValue("all_users", string(dataToCache), 60)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}



func isValidEmail(email string) bool {
    regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(email)
}
