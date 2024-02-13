package transaction

import (
	"baitadores-rinhav2/transaction/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

type Service interface {
	Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, error, int)
	Statement(ctx context.Context, userId int) (dto.StatementResponseDto, error)
}

type serviceImpl struct {
	db *pgxpool.Pool
}

func NewTransaction(db *pgxpool.Pool) Service {
	return &serviceImpl{
		db,
	}
}

func (t serviceImpl) Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, error, int) {

	var result struct {
		Balance string `json:"balance"`
		Limit   string `json:"limit"`
	}

	row := t.db.QueryRow(ctx,
		"SELECT createtransaction($1, $2, $3)",
		userId, validateType(transactionDto.Value, transactionDto.Type), transactionDto.Description,
	)

	if err := row.Scan(&result); err != nil {
		fmt.Println("Erro: resultado não é do tipo esperado")
	}
	var resultado dto.TransactionResponseDto
	if result.Limit == "" {
		errorCode, _ := strconv.Atoi(result.Balance)
		if errorCode == -2 {
			return resultado, fmt.Errorf("saldo insuficiente"), 2
		}
		if errorCode == -1 {
			return resultado, fmt.Errorf("usuario nao encontrado"), 1
		}
		return resultado, fmt.Errorf("Formato de valores inválido"), 0
	}

	saldo, err := strconv.Atoi(result.Balance)
	if err != nil {
		return resultado, fmt.Errorf("Erro ao converter saldo para inteiro"), 0
	}

	limite, err := strconv.Atoi(result.Limit)
	if err != nil {
		return resultado, fmt.Errorf("Erro ao converter limite para inteiro"), 0
	}
	resultado.Balance = saldo
	resultado.Limit = limite * -1

	return resultado, nil, 0
}

func (t serviceImpl) Statement(ctx context.Context, userId int) (dto.StatementResponseDto, error) {
	var response dto.StatementResponseDto
	var resultList []struct {
		Total       int
		Limit       int
		Value       *int
		Description *string
		CreatedAt   *time.Time
	}

	rows, err := t.db.Query(ctx, "SELECT users.balance, users.limit, transaction.value, transaction.description, transaction.created_at "+
		" FROM transaction RIGHT JOIN users on users.id = transaction.user_id "+
		" WHERE users.id = $1 order by transaction.created_at desc limit 10", userId)

	if err != nil {
		panic("error query")
	}

	for rows.Next() {
		var result struct {
			Total       int
			Limit       int
			Value       *int
			Description *string
			CreatedAt   *time.Time
		}
		err = rows.Scan(&result.Total, &result.Limit, &result.Value, &result.Description, &result.CreatedAt)
		if err != nil {
			return dto.StatementResponseDto{}, err
		}
		resultList = append(resultList, result)
	}

	if len(resultList) == 0 {
		return dto.StatementResponseDto{}, fmt.Errorf("User not found")
	}

	response.Balance = dto.BalanceDetails{
		Total:         resultList[0].Total,
		StatementDate: time.Now(),
		Limit:         resultList[0].Limit * -1,
	}

	for _, res := range resultList {
		if res.Value != nil {
			response.LastTransactions = append(response.LastTransactions, dto.LastTransactions{
				Value:       *res.Value * -1,
				Type:        determineType(res.Value),
				Description: *res.Description,
				CreatedAt:   *res.CreatedAt,
			})
		}
	}

	fmt.Printf("%+v\n", resultList)
	return response, nil
}

func validateType(num int, t string) int {
	if t == "c" {
		return num
	} else if t == "d" {
		return -num
	} else {
		return 0
	}
}

func determineType(num *int) string {
	if num != nil && *num >= 0 {
		return "c"
	} else {
		return "d"
	}
}
