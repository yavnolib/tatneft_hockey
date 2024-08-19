package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"tat_hockey_pack/internal/configs/logger"
	"tat_hockey_pack/internal/configs/postgre"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
	l, err := postgre.LoadPgxPool()
	if err != nil {
		panic(err)
	}
	_ = l
	//err = postgre.TestPing(l)
	//if err != nil {
	//	fmt.Println(err)
	//}
	postgre.RemoveTables(l)
	log := logger.InitLogger()

	// Пример логирования
	log.Debug("Debug message", "env", os.Getenv("APP_ENV"))
	log.Info("Application started", "env", os.Getenv("APP_ENV"))
	log.Error("An error occurred", "error", "some error details")

	fmt.Println("Starting server...")
	// TODO: конфиг бд

	// TODO: загрузка конфига логгер

	// TODO middlewares:
	// 1. Request ID
	// 2. Logger
	// 3. Panic Recover
	// 4. CSRF
	// 5. jwt
	// 6. сессии

	// TODO: при получении видео - берем хэш от файла и так сохраняем (дольше, но это сделать все равно надо)
	// либо по шаблону (все равно удаляем)

	// TODO: ручки

	// сервер запуск
}

// TODO к рефакторингу запихнуть основные конфиги в одно место
