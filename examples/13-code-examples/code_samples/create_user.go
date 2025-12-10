// Go example for creating a user
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func main() {
	user := map[string]interface{}{
		"name": "John Doe",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
