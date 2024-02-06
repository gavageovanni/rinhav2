package routes

import (
	"baitadores-rinhav2/config"
	"baitadores-rinhav2/controller"
	"baitadores-rinhav2/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Routing struct {
	nc controller.TransactionController
}

func (r Routing) GetRoutes() *echo.Echo {
	e := echo.New()
	db := config.DB()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	t := transaction.NewTransaction(db)
	nc := controller.NewTransactionController(t)

	e.POST("/clientes/:id/transacoes", nc.Execute)

	return e
}
