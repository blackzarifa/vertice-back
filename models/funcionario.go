package models

import "time"

type Endereco struct {
	ID          int     `json:"id_endereco,omitempty"`
	CEP         string  `json:"cep"`
	Local       string  `json:"local"`
	NumeroCasa  int     `json:"numero_casa"`
	Bairro      string  `json:"bairro"`
	Cidade      string  `json:"cidade"`
	Estado      string  `json:"estado"`
	Complemento *string `json:"complemento,omitempty"`
}

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

type Funcionario struct {
	ID                int     `json:"id_funcionario,omitempty"`
	IDUsuario         int     `json:"id_usuario"`
	IDSupervisor      *int    `json:"id_supervisor,omitempty"`
	CodigoFuncionario string  `json:"codigo_funcionario"`
	Cargo             string  `json:"cargo"`
	Usuario           Usuario `json:"usuario,omitempty"`
}

type CreateFuncionarioRequest struct {
	Nome              string   `json:"nome" binding:"required"`
	CPF               string   `json:"cpf" binding:"required,len=11"`
	DataNascimento    string   `json:"data_nascimento" binding:"required"`
	Telefone          string   `json:"telefone" binding:"required"`
	Senha             string   `json:"senha" binding:"required,min=6"`
	CodigoFuncionario string   `json:"codigo_funcionario" binding:"required"`
	Cargo             string   `json:"cargo" binding:"required,oneof=ESTAGIARIO ATENDENTE GERENTE"`
	IDSupervisor      *int     `json:"id_supervisor,omitempty"`
	Endereco          Endereco `json:"endereco" binding:"required"`
}

type FuncionarioResponse struct {
	ID                int       `json:"id_funcionario"`
	CodigoFuncionario string    `json:"codigo_funcionario"`
	Cargo             string    `json:"cargo"`
	IDSupervisor      *int      `json:"id_supervisor,omitempty"`
	Nome              string    `json:"nome"`
	CPF               string    `json:"cpf"`
	DataNascimento    time.Time `json:"data_nascimento"`
	Telefone          string    `json:"telefone"`
	Endereco          *Endereco `json:"endereco,omitempty"`
}
