package models

type Cliente struct {
	ID           int     `json:"id_cliente,omitempty"`
	IDUsuario    int     `json:"id_usuario"`
	ScoreCredito float64 `json:"score_credito"`
	Usuario      *Usuario `json:"usuario,omitempty"`
}
