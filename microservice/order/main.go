package main

import (
	"encoding/json"
	"microservices/middleware"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Item struct
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Order struct
type Order struct {
	ID    int    `json:"id"`
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

type APIResponse struct {
	Data      interface{} `json:"data"`
	TotalData int         `json:"total_data"`
	Message   string      `json:"message"`
}

// Mock data (in-memory store)
var orders = []Order{
	{
		ID: 1,
		Items: []Item{
			{ID: 1, Name: "Laptop", Price: 5000000},
		},
		Total: 5000000,
	},
	{
		ID: 2,
		Items: []Item{
			{ID: 2, Name: "Mouse", Price: 200000},
			{ID: 3, Name: "Keyboard", Price: 500000},
		},
		Total: 700000,
	},
	{
		ID: 3,
		Items: []Item{
			{ID: 1, Name: "Laptop", Price: 5000000},
			{ID: 2, Name: "Mouse", Price: 200000},
			{ID: 3, Name: "Keyboard", Price: 500000},
		},
		Total: 5700000,
	},
}

// Mutex for thread safety
var mu sync.Mutex

func main() {
	http.HandleFunc("/order", middleware.AuthMiddleware(handleOrders))     // Endpoint untuk GET semua data atau POST
	http.HandleFunc("/order/", middleware.AuthMiddleware(handleOrderByID)) // Endpoint untuk GET, PUT, DELETE berdasarkan ID
	http.ListenAndServe(":8002", nil)
}

// handleOrders handles listing all orders and creating a new order
func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllOrders(w)
	case http.MethodPost:
		createOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleOrderByID handles operations (GET, PUT, DELETE) on individual orders
func handleOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/order/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getOrderByID(w, id)
	case http.MethodPut:
		updateOrder(w, r, id)
	case http.MethodDelete:
		deleteOrder(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllOrders retrieves all orders
func getAllOrders(w http.ResponseWriter) {
	mu.Lock()
	defer mu.Unlock()

	response := APIResponse{
		Data:      orders,
		TotalData: len(orders),
		Message:   "Datas retrieved successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getOrderByID retrieves a single order by its ID
func getOrderByID(w http.ResponseWriter, id int) {
	mu.Lock()
	defer mu.Unlock()

	for _, order := range orders {
		if order.ID == id {
			response := APIResponse{
				Data:      order,
				TotalData: 1,
				Message:   "Data retrieved successfully",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Data not found", http.StatusNotFound)
}

// createOrder adds a new order
func createOrder(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	total := 0
	for _, item := range newOrder.Items {
		total += item.Price
	}
	newOrder.Total = total
	newOrder.ID = len(orders) + 1
	orders = append(orders, newOrder)

	response := APIResponse{
		Data:      newOrder,
		TotalData: len(orders),
		Message:   "Created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// updateOrder updates an existing order
func updateOrder(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	var updatedOrder Order
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for i, order := range orders {
		if order.ID == id {
			updatedOrder.ID = id
			total := 0
			for _, item := range updatedOrder.Items {
				total += item.Price
			}
			updatedOrder.Total = total
			orders[i] = updatedOrder

			response := APIResponse{
				Data:      updatedOrder,
				TotalData: len(orders),
				Message:   "Updated successfully",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Data not found", http.StatusNotFound)
}

// deleteOrder removes an order by its ID
func deleteOrder(w http.ResponseWriter, id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, order := range orders {
		if order.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			response := APIResponse{
				Data:      nil,
				TotalData: len(orders),
				Message:   "Deleted successfully",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Order not found", http.StatusNotFound)
}
