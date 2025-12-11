package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

// CreditRequest - структура вхідних даних
type CreditResponse struct {
	MonthlyPayment float64 `json:"monthly_payment"`
	TotalPayment   float64 `json:"total_payment"`
	Overpayment    float64 `json:"overpayment"`
	Status         string  `json:"status"`
}

// CalculateAnnuity - функція розрахунку (бізнес-логіка)
func CalculateAnnuity(amount float64, rate float64, months int) (float64, float64) {
	if amount <= 0 || rate <= 0 || months <= 0 {
		return 0, 0
	}
	
	monthlyRate := rate / 12 / 100
	
	coefficient := (monthlyRate * math.Pow(1+monthlyRate, float64(months))) / (math.Pow(1+monthlyRate, float64(months)) - 1)
	monthlyPayment := amount * coefficient
	totalPayment := monthlyPayment * float64(months)

	
	monthlyPayment = math.Round(monthlyPayment*100) / 100
	totalPayment = math.Round(totalPayment*100) / 100

	return monthlyPayment, totalPayment
}





// creditHandler - обробляє запити
func creditHandler(w http.ResponseWriter, r *http.Request) {
	
	amountStr := r.URL.Query().Get("amount")
	rateStr := r.URL.Query().Get("rate")
	monthsStr := r.URL.Query().Get("months")

	amount, _ := strconv.ParseFloat(amountStr, 64)
	rate, _ := strconv.ParseFloat(rateStr, 64)
	months, _ := strconv.Atoi(monthsStr)

	monthly, total := CalculateAnnuity(amount, rate, months)

	response := CreditResponse{
		MonthlyPayment: monthly,
		TotalPayment:   total,
		Overpayment:    total - amount,
		Status:         "Calculation Successful",
	}

	if monthly == 0 {
		response.Status = "Invalid Input Data"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/credit", creditHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Credit Calculator Service is Running"})
	})

	port := ":8080"
	fmt.Printf("Banking Service starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}