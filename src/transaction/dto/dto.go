package dto

import "time"

type TransactionRequestDto struct {
	Value       int    `json:"valor" validate:"required"`
	Type        string `json:"tipo" validate:"required"`
	Description string `json:"descricao" validate:"required,max=10"`
}

type TransactionResponseDto struct {
	Balance int `json:"saldo"`
	Limit   int `json:"limite"`
}

type StatementResponseDto struct {
	Balance          BalanceDetails     `json:"saldo"`
	LastTransactions []LastTransactions `json:"ultimas_transacoes"`
}

type BalanceDetails struct {
	Total         int       `json:"total"`
	StatementDate time.Time `json:"data_extrato"`
	Limit         int       `json:"limite"`
}

type LastTransactions struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}
