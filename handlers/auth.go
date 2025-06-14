package handlers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/blackzarifa/vertice-back/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret}
}

// Login authenticates a funcionario and returns a JWT token
// Returns: {"token": "jwt.token.here", "funcionario": {...}}
func (h *AuthHandler) Login(c *gin.Context) {
	// Parse request body
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password to match database storage
	senhaHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Senha)))

	// Find funcionario with matching credentials
	var funcionarioID int
	var nome, codigoFuncionario, cargo string

	query := `
		SELECT f.id_funcionario, u.nome, f.codigo_funcionario, f.cargo
		FROM usuario u
		INNER JOIN funcionario f ON f.id_usuario = u.id_usuario
		WHERE u.cpf = ? AND u.senha_hash = ? AND u.tipo_usuario = 'FUNCIONARIO'
	`
	err := h.db.QueryRow(query, req.CPF, senhaHash).Scan(
		&funcionarioID, &nome, &codigoFuncionario, &cargo,
	)

	// Handle login failures
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Create JWT token with user info
	// Token contains: funcionario_id, nome, cargo
	// Expires in 24 hours
	claims := jwt.MapClaims{
		"funcionario_id": funcionarioID,
		"nome":           nome,
		"cargo":          cargo,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	}

	// Sign token with secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
		return
	}

	// Return token and user info
	response := models.LoginResponse{
		Token: tokenString,
		Funcionario: models.FuncionarioLoginInfo{
			ID:                funcionarioID,
			Nome:              nome,
			CodigoFuncionario: codigoFuncionario,
			Cargo:             cargo,
		},
	}

	c.JSON(http.StatusOK, response)
}
