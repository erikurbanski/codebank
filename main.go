package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/erikurbanski/codebank/infrastructure/grpc/server"
	"github.com/erikurbanski/codebank/infrastructure/kafka"
	"github.com/erikurbanski/codebank/infrastructure/repository"
	"github.com/erikurbanski/codebank/usecase"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello World!")
	db := setupDb()
	defer db.Close()

	// cc := domain.NewCreditCard()
	// cc.Number = "1234"
	// cc.Name = "Erik Urbanski"
	// cc.ExpirationYear = 2025
	// cc.ExpirationMonth = 12
	// cc.CVV = 123
	// cc.Limit = 5380
	// cc.Balance = 0

	// repo := repository.NewTransactionRepositoryDb(db)
	// err := repo.CreateCreditCard(*cc)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	producer := setupKafkaProducer()
	processTransactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(processTransactionUseCase)
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer
	return useCase
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))
	return producer
}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()
	grpcServer.ProcessTransactionUseCase = processTransactionUseCase
	fmt.Println("Rodando gRPC Server")
	grpcServer.Serve()
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
