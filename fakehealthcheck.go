package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.elastic.co/apm/v2"
)

func main() {
	tracer, err := apm.NewTracer("fakehealthcheck", "1.0")
	if err != nil {
		fmt.Println("Error initializing APM agent:", err)
		return
	}
	defer tracer.Close()

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		tx := tracer.StartTransaction("GET", "/api/health")
		defer tx.End()
		tx.Context.SetHTTPRequest(r)

		// Parse the slo parameter
		sloParam := r.URL.Query().Get("slo")
		if sloParam == "" {
			http.Error(w, "Missing 'slo' parameter", http.StatusBadRequest)
			return
		}

		// Convert slo to a float
		slo, err := strconv.ParseFloat(sloParam, 64)
		if err != nil || slo < 0 || slo > 100 {
			http.Error(w, "'slo' must be a decimal number between 0 and 100", http.StatusBadRequest)
			return
		}

		// Seed random number generator
		rand.Seed(time.Now().UnixNano())

		// Generate random boolean based on slo
		randomValue := rand.Float64() * 100 // Random float between 0 and 100
		isSuccessful := randomValue < slo

		// Return response based on random boolean
		if isSuccessful {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "200 OK - Success")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "500 Internal Server Error - Failure")
		}
	})

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
