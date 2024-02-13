package routes

import (
	"baitadores-rinhav2/config2"
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
	config2.Init()
	db := config2.GetDB()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	t := transaction.NewTransaction(db)
	nc := controller.NewTransactionController(t)

	e.POST("/clientes/:id/transacoes", nc.Execute)
	e.GET("/clientes/:id/extrato", nc.Statement)

	return e
}
