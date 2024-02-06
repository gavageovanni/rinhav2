package controller

import (
	"baitadores-rinhav2/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransactionController interface {
	Execute(ctx echo.Context) error
}

type TransactionControllerImpl struct {
	Transaction transaction.Transaction
}

func NewTransactionController(transaction transaction.Transaction) TransactionController {
	return &TransactionControllerImpl{
		transaction,
	}
}

func (t *TransactionControllerImpl) Execute(c echo.Context) error {
	id := c.Param("id")
	var input transaction.TransactionDto
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}

	// Validar o DTO
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "validation failed"})
	}

	userId, _ := strconv.Atoi(id)

	err := t.Transaction.Execute(input, userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, input)
}
