package models

import "time"

type Auditoria struct {
	ID       int       `json:"id_auditoria,omitempty"`
	IDUsuario *int     `json:"id_usuario,omitempty"`
	Acao     string    `json:"acao"`
	DataHora time.Time `json:"data_hora"`
	Detalhes *string   `json:"detalhes,omitempty"`
}
