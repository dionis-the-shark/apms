package main

import (
	"fmt"
	"log"
	"net/http"

	apmsenv "github.com/dionis-the-shark/apms-env"
	"github.com/dionis-the-shark/apms-task-tracker/internal/app"
	internalhttp "github.com/dionis-the-shark/apms-task-tracker/internal/handlers"
	"github.com/dionis-the-shark/apms-task-tracker/internal/storage/postgres"
)

func main() {
	// Збір змінних оточення для підключення до бази даних
	cfg, err := apmsenv.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Підключення до бази даних
	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	// Реєстрація модудлей системи
	services := app.NewServices(db)

	// Реєстрація маршрутів підлючення
	r := internalhttp.NewRouter(services)

	// Запуск сервера
	fmt.Printf("Server started on : %s\n", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r); err != nil {
		panic(fmt.Sprintf("Failed to run server: %v", err))
	}
}
