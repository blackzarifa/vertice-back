package handlers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/blackzarifa/vertice-back/models"
	"github.com/gin-gonic/gin"
)

type FuncionarioHandler struct {
	db *sql.DB
}

func NewFuncionarioHandler(db *sql.DB) *FuncionarioHandler {
	return &FuncionarioHandler{db: db}
}

func (h *FuncionarioHandler) Create(c *gin.Context) {
	var req models.CreateFuncionarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to start transaction"},
		)
		return
	}
	defer tx.Rollback()

	// Insert endereco
	result, err := tx.Exec(
		`
		INSERT INTO endereco (cep, local, numero_casa, bairro, cidade, estado, complemento)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Endereco.CEP,
		req.Endereco.Local,
		req.Endereco.NumeroCasa,
		req.Endereco.Bairro,
		req.Endereco.Cidade,
		req.Endereco.Estado,
		req.Endereco.Complemento,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create address"},
		)
		return
	}
	enderecoID, _ := result.LastInsertId()

	// Parse date and hash password
	dataNascimento, err := time.Parse("2006-01-02", req.DataNascimento)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	senhaHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Senha)))

	// Insert usuario
	result, err = tx.Exec(`
		INSERT INTO usuario (id_endereco, nome, cpf, data_nascimento, telefone, tipo_usuario, senha_hash)
		VALUES (?, ?, ?, ?, ?, 'FUNCIONARIO', ?)`,
		enderecoID, req.Nome, req.CPF, dataNascimento, req.Telefone, senhaHash,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"},
		)
		return
	}
	usuarioID, _ := result.LastInsertId()

	// Insert funcionario
	result, err = tx.Exec(`
		INSERT INTO funcionario (id_usuario, id_supervisor, codigo_funcionario, cargo)
		VALUES (?, ?, ?, ?)`,
		usuarioID, req.IDSupervisor, req.CodigoFuncionario, req.Cargo,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create funcionario"},
		)
		return
	}
	funcionarioID, _ := result.LastInsertId()

	tx.Commit()

	response := models.FuncionarioResponse{
		ID:                int(funcionarioID),
		CodigoFuncionario: req.CodigoFuncionario,
		Cargo:             req.Cargo,
		IDSupervisor:      req.IDSupervisor,
		Nome:              req.Nome,
		CPF:               req.CPF,
		DataNascimento:    dataNascimento,
		Telefone:          req.Telefone,
		Endereco:          req.Endereco,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *FuncionarioHandler) List(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT 
			f.id_funcionario, f.codigo_funcionario, f.cargo, f.id_supervisor,
			u.nome, u.cpf, u.data_nascimento, u.telefone,
			e.cep, e.local, e.numero_casa, e.bairro, e.cidade, e.estado
		FROM funcionario f
		INNER JOIN usuario u ON f.id_usuario = u.id_usuario
		LEFT JOIN endereco e ON u.id_endereco = e.id_endereco
		ORDER BY f.id_funcionario DESC
	`)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch data"},
		)
		return
	}
	defer rows.Close()

	var funcionarios []models.FuncionarioListItem
	for rows.Next() {
		var f models.FuncionarioListItem
		var dataNascimento time.Time
		var idSupervisor sql.NullInt64

		err := rows.Scan(
			&f.ID, &f.CodigoFuncionario, &f.Cargo, &idSupervisor,
			&f.Nome, &f.CPF, &dataNascimento, &f.Telefone,
			&f.CEP, &f.Local, &f.NumeroCasa, &f.Bairro, &f.Cidade, &f.Estado,
		)
		if err != nil {
			continue
		}

		f.DataNascimento = dataNascimento.Format("2006-01-02")
		if idSupervisor.Valid {
			id := int(idSupervisor.Int64)
			f.IDSupervisor = &id
		}

		funcionarios = append(funcionarios, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(funcionarios),
		"data":  funcionarios,
	})
}
