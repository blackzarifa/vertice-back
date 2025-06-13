package models

import "time"

type Usuario struct {
	ID             int       `json:"id_usuario,omitempty"`
	IDEndereco     *int      `json:"id_endereco,omitempty"`
	Nome           string    `json:"nome"`
	CPF            string    `json:"cpf"`
	DataNascimento time.Time `json:"data_nascimento"`
	Telefone       string    `json:"telefone"`
	TipoUsuario    string    `json:"tipo_usuario"`
	SenhaHash      string    `json:"-"`
	Senha          string    `json:"senha,omitempty"`
}
