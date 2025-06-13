package models

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
