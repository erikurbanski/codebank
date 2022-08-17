package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/erikurbanski/codebank/domain"
	"github.com/erikurbanski/codebank/infrastructure/repository"
	"github.com/erikurbanski/codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello World!")
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Erik Urbanski"
	cc.ExpirationYear = 2025
	cc.ExpirationMonth = 12
	cc.CVV = 123
	cc.Limit = 5380
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connection to database!")
	}
	return db
}
