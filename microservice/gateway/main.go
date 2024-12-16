package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, "http://localhost:8001/user/profile")
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, "http://localhost:8002/order")
	})

	http.ListenAndServe(":8000", nil)
}

func proxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
