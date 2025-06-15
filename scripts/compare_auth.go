//go:build ignore
// +build ignore

// Compare password vs OTP authentication
// Run with: go run scripts/compare_auth.go
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

	fmt.Println("=== Authentication Methods Comparison ===\n")

	// Test data
	cpf := "12345678901"
	senha := "senha123"

	// Method 1: Password Authentication
	fmt.Println("METHOD 1: Password Authentication")
	fmt.Println("---------------------------------")
	start := time.Now()

	passwordToken := loginWithPassword(client, cpf, senha)

	fmt.Printf("Time taken: %v\n", time.Since(start))
	fmt.Printf("Steps: 1 (send CPF + password)\n")
	fmt.Printf("User experience: Simple, traditional\n")

	// Method 2: OTP Authentication
	fmt.Println("\n\nMETHOD 2: OTP Authentication")
	fmt.Println("----------------------------")
	start = time.Now()

	// Request OTP
	otp := requestOTP(client, cpf)
	if otp != "" {
		// Verify OTP
		otpToken := verifyOTP(client, cpf, otp)

		fmt.Printf("Time taken: %v\n", time.Since(start))
		fmt.Printf("Steps: 2 (request OTP, then verify)\n")
		fmt.Printf("User experience: More secure, modern\n")

		if otpToken != "" && passwordToken != "" {
			fmt.Println("\n✓ Both methods working correctly!")
		}
	}
}

func loginWithPassword(client *http.Client, cpf, senha string) string {
	reqBody := map[string]string{
		"cpf":   cpf,
		"senha": senha,
	}
	jsonData, _ := json.Marshal(reqBody)

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

	if resp.StatusCode == 200 {
		var result map[string]any
		json.Unmarshal(body, &result)
		if token, ok := result["token"].(string); ok {
			fmt.Printf("✓ Password login successful\n")
			return token
		}
	} else {
		fmt.Printf("✗ Password login failed: %s\n", body)
	}

	return ""
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

	if resp.StatusCode == 200 {
		var result map[string]any
		json.Unmarshal(body, &result)
		if otp, ok := result["otp"].(string); ok {
			fmt.Printf("✓ OTP requested: %s\n", otp)
			return otp
		}
	} else {
		fmt.Printf("✗ OTP request failed: %s\n", body)
	}

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

	if resp.StatusCode == 200 {
		var result map[string]any
		json.Unmarshal(body, &result)
		if token, ok := result["token"].(string); ok {
			fmt.Printf("✓ OTP verified successfully\n")
			return token
		}
	} else {
		fmt.Printf("✗ OTP verification failed: %s\n", body)
	}

	return ""
}
