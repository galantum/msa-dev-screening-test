package main

import (
	"encoding/json"
	"microservices/middleware"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password,omitempty"`
}

type APIResponse struct {
	Data      interface{} `json:"data"`
	TotalData int         `json:"total_data"`
	Message   string      `json:"message"`
}

// Mock data (in-memory store)
var users = []User{
	{ID: 1, Username: "manager", Email: "manager@example.com", Age: 40, Password: "manager"},
	{ID: 2, Username: "staff", Email: "staff@example.com", Age: 21, Password: "staff"},
	{ID: 3, Username: "admin", Email: "admin@example.com", Age: 25, Password: "admin"},
}

var (
	secretKey = "msa_grandis_2024" // Ganti dengan secret key yang lebih aman
	mu        sync.Mutex           // Mutex for thread safety
)

func main() {
	http.HandleFunc("/user/login", Login)
	http.HandleFunc("/user/profile/", middleware.AuthMiddleware(HandleProfileRequests)) // GET/PUT specific profile
	http.ListenAndServe(":8001", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate email & password
	for _, user := range users {
		if user.Email == credentials.Email && user.Password == credentials.Password {
			// Buat token JWT
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userID": user.ID,
				"exp":    time.Now().Add(time.Hour * 1).Unix(), // Token expired is 1 hours
			})

			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
				http.Error(w, "Could not generate token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
			return
		}
	}

	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
}

// HandleProfileRequests - Handle GET and PUT requests for individual profiles
func HandleProfileRequests(w http.ResponseWriter, r *http.Request) {
	userName := strings.TrimPrefix(r.URL.Path, "/user/profile/")
	if userName == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetProfile(w, r, userName)
	case http.MethodPut:
		UpdateProfile(w, r, userName)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetProfile - Get user profile by username
func GetProfile(w http.ResponseWriter, r *http.Request, userName string) {
	mu.Lock()
	defer mu.Unlock()

	for _, user := range users {
		if strings.EqualFold(user.Username, userName) {
			user.Password = ""

			response := APIResponse{
				Data:      user,
				TotalData: 1, // Only 1 data found
				Message:   "Data retrieved successfully",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// If user not found
	response := APIResponse{
		Data:      nil,
		TotalData: 0,
		Message:   "Data not found",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile - Update user profile
func UpdateProfile(w http.ResponseWriter, r *http.Request, userName string) {
	mu.Lock()
	defer mu.Unlock()

	// Parse the body
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if strings.EqualFold(user.Username, userName) {
			// Update the user data
			users[i].Username = updatedUser.Username
			users[i].Email = updatedUser.Email
			users[i].Age = updatedUser.Age

			response := APIResponse{
				Data:      users[i],
				TotalData: 1, // Only 1 data found
				Message:   "Updated successfully",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := APIResponse{
		Data:      nil,
		TotalData: 0,
		Message:   "Not found data!",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}
