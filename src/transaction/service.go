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
	Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, *int)
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

type ExecuteResult struct {
	Balance string `json:"balance"`
	Limit   string `json:"limit"`
}

func (t serviceImpl) Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, *int) {
	var result ExecuteResult
	row := t.db.QueryRow(ctx,
		"SELECT createtransaction($1, $2, $3)",
		userId, validateType(transactionDto.Value, transactionDto.Type), transactionDto.Description,
	)
	if err := row.Scan(&result); err != nil {
		fmt.Println("deserialize sql result error")
	}
	return buildExecuteResponse(result)
}

func buildExecuteResponse(result ExecuteResult) (dto.TransactionResponseDto, *int) {
	var transactionResponse dto.TransactionResponseDto
	var errorCode *int
	balance, err := strconv.Atoi(result.Balance)
	if err != nil {
		return transactionResponse, nil
	}

	if result.Limit == "" {
		errorCode = &balance
		return transactionResponse, errorCode
	}

	limit, err := strconv.Atoi(result.Limit)
	if err != nil {
		unknownErrorCode := -3
		errorCode = &unknownErrorCode
		return transactionResponse, errorCode
	}
	transactionResponse.Balance = balance
	transactionResponse.Limit = limit * -1

	return transactionResponse, nil
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
				Value:       makePositive(res.Value),
				Type:        determineType(res.Value),
				Description: *res.Description,
				CreatedAt:   *res.CreatedAt,
			})
		}
	}
	return response, nil
}

func makePositive(numberPtr *int) int {
	if numberPtr != nil {
		num := *numberPtr
		if num < 0 {
			return -num
		}
		return num
	}
	return 0
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
