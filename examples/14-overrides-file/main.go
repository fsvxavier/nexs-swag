package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// @title Overrides File API
// @version 1.0
// @host localhost:8080
// @BasePath /api

// Account demonstrates type overrides.
type Account struct {
	ID        sql.NullInt64  `json:"id"`         // Will be overridden to integer
	Balance   sql.NullString `json:"balance"`    // Will be overridden to string
	CreatedAt time.Time      `json:"created_at"` // Will keep time format
	UpdatedAt *time.Time     `json:"updated_at"` // Optional time
}

// GetAccount returns an account
// @Summary Get account
// @Tags accounts
// @Produce json
// @Success 200 {object} Account
// @Router /accounts/{id} [get].
func GetAccount(w http.ResponseWriter, r *http.Request) {
	account := Account{
		ID:        sql.NullInt64{Int64: 123, Valid: true},
		Balance:   sql.NullString{String: "1000.50", Valid: true},
		CreatedAt: time.Now(),
	}
	json.NewEncoder(w).Encode(account)
}

func main() {
	http.HandleFunc("/api/accounts/", GetAccount)
	http.ListenAndServe(":8080", nil)
}
