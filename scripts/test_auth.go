//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	// Test login
	fmt.Println("\nTesting login...")
	token := testLogin(client)

	if token != "" {
		// Test protected endpoints
		fmt.Println("\nTesting protected endpoints with token...")
		testProtectedEndpoints(client, token)

		// Test without token
		fmt.Println("\nTesting protected endpoints WITHOUT token...")
		testProtectedEndpoints(client, "")
	}
}

func testLogin(client *http.Client) string {
	// Try to login with a funcionario (you'll need to create one first)
	loginData := map[string]string{
		"cpf":   "12345678901", // Change this to a real CPF from your database
		"senha": "senha123",    // Change this to match
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := client.Post(
		baseURL+"/auth/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Login response: %d - %s\n", resp.StatusCode, body)

	if resp.StatusCode == 200 {
		var loginResp map[string]any
		json.Unmarshal(body, &loginResp)
		if token, ok := loginResp["token"].(string); ok {
			fmt.Printf("Got token: %s...\n", token[:20])
			return token
		}
	}
	return ""
}

func testProtectedEndpoints(client *http.Client, token string) {
	// Test GET /funcionarios
	req, _ := http.NewRequest("GET", baseURL+"/funcionarios", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf(
		"GET /funcionarios: %d - %s\n",
		resp.StatusCode,
		string(body)[:100]+"...",
	)
}
