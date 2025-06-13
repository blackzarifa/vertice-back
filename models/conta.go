package models

import "time"

type Conta struct {
	ID               int       `json:"id_conta,omitempty"`
	NumeroConta      string    `json:"numero_conta"`
	DigitoVerificador string   `json:"digito_verificador"`
	IDAgencia        int       `json:"id_agencia"`
	Saldo            float64   `json:"saldo"`
	TipoConta        string    `json:"tipo_conta"`
	IDCliente        int       `json:"id_cliente"`
	DataAbertura     time.Time `json:"data_abertura"`
	Status           string    `json:"status"`
}

type ContaPoupanca struct {
	ID               int       `json:"id_conta_poupanca,omitempty"`
	IDConta          int       `json:"id_conta"`
	TaxaRendimento   float64   `json:"taxa_rendimento"`
	UltimoRendimento *time.Time `json:"ultimo_rendimento,omitempty"`
}

type ContaCorrente struct {
	ID              int       `json:"id_conta_corrente,omitempty"`
	IDConta         int       `json:"id_conta"`
	Limite          float64   `json:"limite"`
	DataVencimento  time.Time `json:"data_vencimento"`
	TaxaManutencao  float64   `json:"taxa_manutencao"`
}

type ContaInvestimento struct {
	ID                 int     `json:"id_conta_investimento,omitempty"`
	IDConta            int     `json:"id_conta"`
	PerfilRisco        string  `json:"perfil_risco"`
	ValorMinimo        float64 `json:"valor_minimo"`
	TaxaRendimentoBase float64 `json:"taxa_rendimento_base"`
}
