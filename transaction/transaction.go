package transaction

import (
	"baitadores-rinhav2/dto"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"regexp"
	"strconv"
	"time"
)

type Transaction interface {
	Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, error, int)
	Statement(ctx context.Context, userId int) (dto.StatementResponseDto, error)
}

type transactionImpl struct {
	db *pgxpool.Pool
}

func NewTransaction(db *pgxpool.Pool) Transaction {
	return &transactionImpl{
		db,
	}
}

func (t transactionImpl) Execute(ctx context.Context, transactionDto dto.TransactionRequestDto, userId int) (dto.TransactionResponseDto, error, int) {

	//var result struct {
	//	Createtransaction string
	//}

	//if err := t.db.Raw("SELECT createtransaction($1, $2, $3 )",
	//	userId, validateType(transactionDto.Value, transactionDto.Type),
	//	transactionDto.Description,
	//).Scan(&result).Error; err != nil {
	//	panic("falha ao chamar a função armazenada")
	//}

	//tx, err := t.db.Begin(ctx)
	//if err != nil {
	//	panic("falha ao conectar DB")
	//}
	//defer tx.Rollback(ctx)
	//
	//_, err = tx.Exec(ctx,
	//	"SELECT createtransaction($1, $2, $3 ) VALUES ($1, $2, $3)",
	//	userId, validateType(transactionDto.Value, transactionDto.Type), transactionDto.Description,
	//)
	var resultado any // Usando []byte para lidar com o resultado

	row := t.db.QueryRow(ctx,
		"SELECT createtransaction($1, $2, $3)",
		userId, validateType(transactionDto.Value, transactionDto.Type), transactionDto.Description,
	)

	if err := row.Scan(&resultado); err != nil {
		println(err)
	}

	response, err, errorCode := parseTransactionString("-2")
	if err != nil {
		fmt.Println("Erro ao fazer o parse da string:", err)
		return dto.TransactionResponseDto{}, err, errorCode
	}

	fmt.Printf("Resultado: %+v\n", response)

	return response, nil, 0
}

func (t transactionImpl) Statement(ctx context.Context, userId int) (dto.StatementResponseDto, error) {
	var response dto.StatementResponseDto
	var result []struct {
		Total       int       `gorm:"column:balance"`
		Limit       int       `gorm:"column:limit"`
		Value       int       `gorm:"column:value"`
		Description string    `gorm:"column:description"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	//t.db.Table("user_models").
	//	Select("user_models.balance, user_models.limit, models.value, models.description, models.created_at").
	//	Joins("LEFT JOIN models ON user_models.id = models.user_id").
	//	Where("user_models.id = ?", userId).
	//	Order("models.created_at desc").
	//	Limit(10).
	//	Find(&result)

	err := t.db.QueryRow(ctx, "SELECT user_models.balance, user_models.limit, models.value, models.description, models.created_at "+
		" FROM models on user_models.id = models.user_id "+
		" WHERE public.clientes WHERE id = $1 order by models.created_at desc ",
		userId).Scan(&result)
	if err != nil {
		panic("error query")
	}

	if len(result) == 0 {
		return dto.StatementResponseDto{}, fmt.Errorf("User not found")
	}

	response.Balance = dto.BalanceDetails{
		Total:         result[0].Total,
		StatementDate: time.Now(),
		Limit:         result[0].Limit * -1,
	}

	for _, res := range result {
		response.LastTransactions = append(response.LastTransactions, dto.LastTransactions{
			Value:       res.Value,
			Type:        determineType(res.Value),
			Description: res.Description,
			CreatedAt:   res.CreatedAt,
		})
	}

	fmt.Printf("%+v\n", result)
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

func parseTransactionString(input string) (dto.TransactionResponseDto, error, int) {
	var resultado dto.TransactionResponseDto

	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 || len(matches[0]) != 2 {
		return resultado, fmt.Errorf("String inválida"), 0
	}

	values := matches[0][1]

	splitValues := regexp.MustCompile(`,`).Split(values, 2)

	if len(splitValues) != 2 {
		errorCode, _ := strconv.Atoi(splitValues[0])
		if errorCode == -2 {
			return resultado, fmt.Errorf("saldo insuficiente"), 2
		}
		if errorCode == -1 {
			return resultado, fmt.Errorf("usuario nao encontrado"), 1
		}

		return resultado, fmt.Errorf("Formato de valores inválido"), 0
	}

	saldo, err := strconv.Atoi(splitValues[0])
	if err != nil {
		return resultado, fmt.Errorf("Erro ao converter saldo para inteiro"), 0
	}

	limite, err := strconv.Atoi(splitValues[1])
	if err != nil {
		return resultado, fmt.Errorf("Erro ao converter limite para inteiro"), 0
	}

	resultado.Balance = saldo
	resultado.Limit = limite * -1

	return resultado, nil, 0
}

func determineType(num int) string {
	if num >= 0 {
		return "c"
	} else {
		return "d"
	}
}
