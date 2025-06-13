package models

type Agencia struct {
	ID            int       `json:"id_agencia,omitempty"`
	Nome          string    `json:"nome"`
	CodigoAgencia string    `json:"codigo_agencia"`
	EnderecoID    *int      `json:"endereco_id,omitempty"`
	Endereco      *Endereco `json:"endereco,omitempty"`
}
