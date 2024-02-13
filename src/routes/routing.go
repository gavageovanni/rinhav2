package routes

import (
	"baitadores-rinhav2/config"
	"baitadores-rinhav2/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Routing struct {
	nc transaction.Controller
}

func (r Routing) GetRoutes() *echo.Echo {
	e := echo.New()
	config.Init()
	db := config.GetDB()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	t := transaction.NewTransaction(db)
	nc := transaction.NewTransactionController(t)

	e.POST("/clientes/:id/transacoes", nc.Execute)
	e.GET("/clientes/:id/extrato", nc.Statement)

	return e
}
