package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	ID      int     `gorm:"primaryKey;autoIncrement"`
	Balance float64 `gorm:"type:decimal(10,2);not null;default:0"`
}
type Transaction struct {
	ID            int     `gorm:"primaryKey;autoIncrement"`
	FromAccountID uint    `gorm:"not null"`
	ToAccountID   uint    `gorm:"not null"`
	Amount        float64 `gorm:"type:decimal(10,2);not null;default:0"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	db.AutoMigrate(&Account{}, &Transaction{})
	initAccount(db)

}
func initAccount(db *gorm.DB) {
	account := []Account{
		{Balance: 1000},
		{Balance: 2000},
	}
	db.Create(&account)
}
func createTransaction(db *gorm.DB, fromId, toId uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var fromAccount Account
		var toAccount Account
		if error := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromAccount, fromId).Error; error != nil {
			return fmt.Errorf("获取账户失败：%v", error)
		}
		if fromAccount.Balance < amount {
			return errors.New("余额不足")
		}
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
			return fmt.Errorf("扣款失败：%v", err)
		}
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
			return fmt.Errorf("收款失败：%v", err)
		}
		transaction := Transaction{
			FromAccountID: fromId,
			ToAccountID:   toId,
			Amount:        amount,
		}
		if error := tx.Create(&transaction).Error; error != nil {
			return fmt.Errorf("创建交易记录失败：%v", error)
		}
		return nil
	})
}
