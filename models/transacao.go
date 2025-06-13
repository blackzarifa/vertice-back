package models

import "time"

type Transacao struct {
	ID             int       `json:"id_transacao,omitempty"`
	IDContaOrigem  *int      `json:"id_conta_origem,omitempty"`
	IDContaDestino *int      `json:"id_conta_destino,omitempty"`
	TipoTransacao  string    `json:"tipo_transacao"`
	Valor          float64   `json:"valor"`
	DataHora       time.Time `json:"data_hora"`
	Descricao      *string   `json:"descricao,omitempty"`
}
