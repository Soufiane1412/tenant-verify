package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/soufiane1412/tenant-verify/internal/verification"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Tenant Verify starting on port %s\n", port)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/", rootHandler)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// Only accepting POST requests
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}
	var tenant verification.TenantRequest
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	// run verification
	result, err := verification.VerifyTenant(tenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log and return result
	log.Printf("Verified: %s, Status: %s", result.ID, result.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

// Decode JSON format req

func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"health","service":"tenant-verify"}`))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Tenant Verify API"))
}
