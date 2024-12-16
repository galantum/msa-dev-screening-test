package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type APIResponse struct {
	Data      interface{} `json:"data"`
	TotalData int         `json:"total_data"`
	Message   string      `json:"message"`
}

// Mock data (in-memory store)
var users = []User{
	{ID: 1, Username: "manager", Email: "manager@example.com", Age: 40},
	{ID: 2, Username: "staff", Email: "staff@example.com", Age: 21},
	{ID: 3, Username: "admin", Email: "admin@example.com", Age: 25},
}

// Mutex for thread safety
var mu sync.Mutex

func main() {
	http.HandleFunc("/user", GetAllProfiles)                 // GET all profiles
	http.HandleFunc("/user/profile/", HandleProfileRequests) // GET/PUT specific profile
	http.ListenAndServe(":8001", nil)
}

// GetAllProfiles - Get all user profiles
func GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := "No data!"

	if len(users) > 1 {
		msg = "Get success!"
	}

	// Buat response sesuai format
	response := APIResponse{
		Data:      users,
		TotalData: len(users), // Total user data
		Message:   msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
			// Bungkus response
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

			// Bungkus response
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
