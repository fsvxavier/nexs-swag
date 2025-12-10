package main

import (
	"encoding/json"
	"net/http"
)

// @title Collection Format API
// @version 1.0
// @host localhost:8080
// @BasePath /api

type SearchRequest struct {
	Tags []string `json:"tags"`
}

// SearchItems searches items by tags
// @Summary Search items
// @Tags search
// @Accept json
// @Produce json
// @Param tags query []string true "Tags to search" collectionFormat(multi)
// @Success 200 {array} string
// @Router /search [get]
func SearchItems(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tags"]
	json.NewEncoder(w).Encode(tags)
}

// FilterItems filters with CSV format
// @Summary Filter items
// @Tags search
// @Param ids query []int true "IDs" collectionFormat(csv)
// @Success 200 {array} int
// @Router /filter [get]
func FilterItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]int{1, 2, 3})
}

// PipeItems demonstrates pipe format
// @Summary Pipe separated
// @Tags search
// @Param statuses query []string true "Statuses" collectionFormat(pipes)
// @Success 200 {array} string
// @Router /pipe [get]
func PipeItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"active", "pending"})
}

func main() {
	http.HandleFunc("/api/search", SearchItems)
	http.HandleFunc("/api/filter", FilterItems)
	http.HandleFunc("/api/pipe", PipeItems)
	http.ListenAndServe(":8080", nil)
}
