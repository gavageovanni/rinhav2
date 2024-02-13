package routes

import (
	"baitadores-rinhav2/config"
	"baitadores-rinhav2/transaction"
	"github.com/labstack/echo/v4"
)

type Routing struct {
	nc transaction.Controller
}

func (r Routing) GetRoutes() *echo.Echo {
	e := echo.New()
	err := config.Init()
	if err != nil {
		return nil
	}
	db := config.GetDB()

	t := transaction.NewTransaction(db)
	nc := transaction.NewTransactionController(t)

	e.POST("/clientes/:id/transacoes", nc.Execute)
	e.GET("/clientes/:id/extrato", nc.Statement)

	return e
}
