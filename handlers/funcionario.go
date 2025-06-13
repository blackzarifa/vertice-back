package handlers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blackzarifa/vertice-back/models"
	"github.com/gin-gonic/gin"
)

func CreateFuncionario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateFuncionarioRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()

		var enderecoID int
		result, err := tx.Exec(`
			INSERT INTO endereco (cep, local, numero_casa, bairro, cidade, estado, complemento)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			req.Endereco.CEP, req.Endereco.Local, req.Endereco.NumeroCasa,
			req.Endereco.Bairro, req.Endereco.Cidade, req.Endereco.Estado,
			req.Endereco.Complemento)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
			return
		}
		
		enderecoID64, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get address ID"})
			return
		}
		enderecoID = int(enderecoID64)

		dataNascimento, err := time.Parse("2006-01-02", req.DataNascimento)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}

		senhaHash := fmt.Sprintf("%x", md5.Sum([]byte(req.Senha)))

		var usuarioID int
		result, err = tx.Exec(`
			INSERT INTO usuario (id_endereco, nome, cpf, data_nascimento, telefone, tipo_usuario, senha_hash)
			VALUES (?, ?, ?, ?, ?, 'FUNCIONARIO', ?)`,
			enderecoID, req.Nome, req.CPF, dataNascimento, req.Telefone, senhaHash)
		
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "cpf") {
				c.JSON(http.StatusConflict, gin.H{"error": "CPF already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			}
			return
		}
		
		usuarioID64, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}
		usuarioID = int(usuarioID64)

		var funcionarioID int
		result, err = tx.Exec(`
			INSERT INTO funcionario (id_usuario, id_supervisor, codigo_funcionario, cargo)
			VALUES (?, ?, ?, ?)`,
			usuarioID, req.IDSupervisor, req.CodigoFuncionario, req.Cargo)
		
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "codigo_funcionario") {
				c.JSON(http.StatusConflict, gin.H{"error": "Codigo funcionario already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create funcionario: " + err.Error()})
			}
			return
		}
		
		funcionarioID64, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get funcionario ID"})
			return
		}
		funcionarioID = int(funcionarioID64)

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		response := models.FuncionarioResponse{
			ID:                funcionarioID,
			CodigoFuncionario: req.CodigoFuncionario,
			Cargo:             req.Cargo,
			IDSupervisor:      req.IDSupervisor,
			Nome:              req.Nome,
			CPF:               req.CPF,
			DataNascimento:    dataNascimento,
			Telefone:          req.Telefone,
			Endereco:          req.Endereco,
		}
		response.Endereco.ID = enderecoID

		c.JSON(http.StatusCreated, response)
	}
}

func ListFuncionarios(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT 
				f.id_funcionario,
				f.codigo_funcionario,
				f.cargo,
				f.id_supervisor,
				u.id_usuario,
				u.nome,
				u.cpf,
				u.data_nascimento,
				u.telefone,
				e.id_endereco,
				e.cep,
				e.local,
				e.numero_casa,
				e.bairro,
				e.cidade,
				e.estado,
				e.complemento
			FROM funcionario f
			INNER JOIN usuario u ON f.id_usuario = u.id_usuario
			LEFT JOIN endereco e ON u.id_endereco = e.id_endereco
			ORDER BY f.id_funcionario DESC`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch funcionarios"})
			return
		}
		defer rows.Close()

		var funcionarios []models.FuncionarioResponse
		for rows.Next() {
			var f models.FuncionarioResponse
			var e models.Endereco
			var enderecoID, usuarioID sql.NullInt64
			var complemento sql.NullString

			err := rows.Scan(
				&f.ID,
				&f.CodigoFuncionario,
				&f.Cargo,
				&f.IDSupervisor,
				&usuarioID,
				&f.Nome,
				&f.CPF,
				&f.DataNascimento,
				&f.Telefone,
				&enderecoID,
				&e.CEP,
				&e.Local,
				&e.NumeroCasa,
				&e.Bairro,
				&e.Cidade,
				&e.Estado,
				&complemento,
			)
			if err != nil {
				continue
			}

			if enderecoID.Valid {
				e.ID = int(enderecoID.Int64)
				if complemento.Valid {
					e.Complemento = &complemento.String
				}
				f.Endereco = &e
			}

			funcionarios = append(funcionarios, f)
		}

		if funcionarios == nil {
			funcionarios = []models.FuncionarioResponse{}
		}

		c.JSON(http.StatusOK, gin.H{
			"count": len(funcionarios),
			"data":  funcionarios,
		})
	}
}
