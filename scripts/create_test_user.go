//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Create a test funcionario
	funcionarioData := map[string]interface{}{
		"nome":               "João Silva",
		"cpf":                "12345678901",
		"data_nascimento":    "1990-05-15",
		"telefone":           "11987654321",
		"senha":              "senha123",
		"codigo_funcionario": "FUNC001",
		"cargo":              "GERENTE",
		"endereco": map[string]interface{}{
			"cep":         "01310-100",
			"local":       "Avenida Paulista",
			"numero_casa": 1578,
			"bairro":      "Bela Vista",
			"cidade":      "São Paulo",
			"estado":      "SP",
			"complemento": "Apto 101",
		},
	}

	jsonData, _ := json.Marshal(funcionarioData)
	resp, err := http.Post(
		"http://localhost:8080/api/v1/funcionarios",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error creating funcionario: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Create funcionario response: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", body)

	if resp.StatusCode == 201 {
		fmt.Println("\nFuncionario created successfully!")
		fmt.Println("You can now login with:")
		fmt.Println("CPF: 12345678901")
		fmt.Println("Senha: senha123")
	}
}
