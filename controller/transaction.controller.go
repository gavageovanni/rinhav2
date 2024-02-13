package controller

import (
	"baitadores-rinhav2/dto"
	"baitadores-rinhav2/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransactionController interface {
	Execute(ctx echo.Context) error
	Statement(ctx echo.Context) error
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
	var input dto.TransactionRequestDto
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid format"})
	}

	if input.Type != "c" && input.Type != "d" {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid format"})
	}

	// Validar o DTO
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	userId, _ := strconv.Atoi(id)

	resp, err, errorCode := t.Transaction.Execute(c.Request().Context(), input, userId)
	if err != nil {
		if errorCode == 2 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "insuficient balance"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, resp)
}

func (t *TransactionControllerImpl) Statement(c echo.Context) error {
	id := c.Param("id")

	userId, _ := strconv.Atoi(id)

	resp, err := t.Transaction.Statement(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, resp)
	}

	return c.JSON(http.StatusOK, resp)
}
