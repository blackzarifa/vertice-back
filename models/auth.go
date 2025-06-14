package models

type LoginRequest struct {
	CPF   string `json:"cpf" binding:"required,len=11"`
	Senha string `json:"senha" binding:"required"`
}

type LoginResponse struct {
	Token       string                 `json:"token"`
	Funcionario FuncionarioLoginInfo `json:"funcionario"`
}

type FuncionarioLoginInfo struct {
	ID                int    `json:"id_funcionario"`
	Nome              string `json:"nome"`
	CodigoFuncionario string `json:"codigo_funcionario"`
	Cargo             string `json:"cargo"`
}
