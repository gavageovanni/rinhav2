package transaction

import (
	"gorm.io/gorm"
	"time"
)

type Transaction interface {
	Execute(transactionDto TransactionDto, userId int) error
}

type transactionImpl struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) Transaction {
	return &transactionImpl{
		db,
	}
}

func (t transactionImpl) Execute(transactionDto TransactionDto, userId int) error {

	if err := t.db.Create(&TransactionModel{
		Value:       transactionDto.Value,
		UserId:      userId,
		Description: transactionDto.Description,
		Type:        transactionDto.Type,
	}).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		println(data)
	}

	return nil
}

type UserModel struct {
	gorm.Model
	Id      int    `gorm:"type:int;primary_key"`
	Name    string `gorm:"type:varchar(255)"`
	Limit   int    `gorm:"type:int"`
	Balance int    `gorm:"type:int"`
}

type TransactionModel struct {
	gorm.Model
	Id          int       `gorm:"type:int;primary_key"`
	Value       int       `gorm:"type:int"`
	Type        string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
	UserId      int       `gorm:"type:int"`
}

type TransactionDto struct {
	Value       int    `json:"value" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description" validate:"required"`
}
