package models

import "time"

type Relatorio struct {
	ID            int       `json:"id_relatorio,omitempty"`
	IDFuncionario int       `json:"id_funcionario"`
	TipoRelatorio string    `json:"tipo_relatorio"`
	DataGeracao   time.Time `json:"data_geracao"`
	Conteudo      string    `json:"conteudo"`
}
