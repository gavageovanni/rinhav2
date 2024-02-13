package transaction

import (
	"baitadores-rinhav2/transaction/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Controller interface {
	Execute(ctx echo.Context) error
	Statement(ctx echo.Context) error
}

type ControllerImpl struct {
	Transaction Service
}

func NewTransactionController(service Service) Controller {
	return &ControllerImpl{
		service,
	}
}

func (t *ControllerImpl) Execute(c echo.Context) error {
	id := c.Param("id")
	var input dto.TransactionRequestDto
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid format"})
	}

	if input.Type != "c" && input.Type != "d" {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "invalid type"})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	userId, _ := strconv.Atoi(id)

	resp, errorCode := t.Transaction.Execute(c.Request().Context(), input, userId)
	if errorCode != nil {
		if *errorCode == -2 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "insuficient balance"})
		}
		if *errorCode == -1 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusNotFound, map[string]string{})
	}

	return c.JSON(http.StatusOK, resp)
}

func (t *ControllerImpl) Statement(c echo.Context) error {
	id := c.Param("id")

	userId, _ := strconv.Atoi(id)

	resp, err := t.Transaction.Statement(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, resp)
	}

	return c.JSON(http.StatusOK, resp)
}
