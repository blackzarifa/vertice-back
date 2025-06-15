//go:build ignore
// +build ignore

// Test script for OTP authentication flow
// Run with: go run scripts/test_otp.go
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

	// Test with an existing user
	cpf := "12345678901" // João Silva from create_test_user.go

	fmt.Println("=== OTP Authentication Test ===")
	fmt.Printf("Testing with CPF: %s\n\n", cpf)

	// Step 1: Request OTP
	fmt.Println("1. Requesting OTP...")
	otp := requestOTP(client, cpf)

	if otp == "" {
		fmt.Println("Failed to get OTP. Make sure user exists.")
		return
	}

	// Step 2: Try with wrong OTP
	fmt.Println("\n2. Testing with WRONG OTP...")
	verifyOTP(client, cpf, "000000")

	// Step 3: Try with correct OTP
	fmt.Println("\n3. Testing with CORRECT OTP...")
	token := verifyOTP(client, cpf, otp)

	// Step 4: Use the token
	if token != "" {
		fmt.Println("\n4. Testing protected endpoint with OTP token...")
		testProtectedEndpoint(client, token)
	}

	// Step 5: Try to reuse the same OTP
	fmt.Println("\n5. Testing OTP reuse (should fail)...")
	verifyOTP(client, cpf, otp)
}

func requestOTP(client *http.Client, cpf string) string {
	reqBody := map[string]string{"cpf": cpf}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := client.Post(
		baseURL+"/auth/otp/request",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]any
	json.Unmarshal(body, &result)

	fmt.Printf("Response: %d\n", resp.StatusCode)
	if otp, ok := result["otp"].(string); ok {
		fmt.Printf("OTP received: %s\n", otp)
		fmt.Printf("User: %s\n", result["user_name"])
		fmt.Printf("Expires in: %s\n", result["expires_in"])
		return otp
	}

	fmt.Printf("Error: %s\n", body)
	return ""
}

func verifyOTP(client *http.Client, cpf, otp string) string {
	reqBody := map[string]string{
		"cpf": cpf,
		"otp": otp,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := client.Post(
		baseURL+"/auth/otp/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]any
	json.Unmarshal(body, &result)

	fmt.Printf("Response: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		if token, ok := result["token"].(string); ok {
			fmt.Printf("Login successful!\n")
			fmt.Printf("Token: %s...\n", token[:20])
			if funcionario, ok := result["funcionario"].(map[string]any); ok {
				fmt.Printf("Funcionario: %s (ID: %.0f, Cargo: %s)\n",
					funcionario["nome"],
					funcionario["id_funcionario"],
					funcionario["cargo"])
			}
			return token
		}
	} else {
		fmt.Printf("Error: %s\n", result["error"])
	}

	return ""
}

func testProtectedEndpoint(client *http.Client, token string) {
	req, _ := http.NewRequest("GET", baseURL+"/funcionarios", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Protected endpoint response: %d\n", resp.StatusCode)
	if resp.StatusCode == 200 {
		fmt.Println("✓ Token is valid and working!")
	}
}
