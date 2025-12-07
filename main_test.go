package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Тестуємо чисту математику
func TestCalculateAnnuity(t *testing.T) {
	
	payment, total := CalculateAnnuity(10000, 12, 12)

	if payment != 888.49 {
		t.Errorf("Expected 888.49, got %f", payment)
	}

	if total <= 10000 {
		t.Errorf("Total payment should be greater than amount")
	}
}

//Тестуємо захист (негативні числа)
func TestCalculateZero(t *testing.T) {
	payment, _ := CalculateAnnuity(0, 10, 10)
	if payment != 0 {
		t.Error("Expected 0 payment for 0 amount")
	}
}

// Тестуємо сам веб-сервер
func TestCreditHandler(t *testing.T) {
	
	req, err := http.NewRequest("GET", "/credit?amount=10000&rate=12&months=12", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(creditHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response CreditResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.MonthlyPayment != 888.49 {
		t.Errorf("Expected monthly payment 888.49, got %v", response.MonthlyPayment)
	}
}