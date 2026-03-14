package main

import (
	"log"
	"log/slog"
)

func main() {
	// ТЕСТ ПРАВИЛА 1: Строчная буква
	slog.Info("Starting server on port 8080")   // ДОЛЖНО быть ошибкой
	slog.Error("Failed to connect to database") // ДОЛЖНО быть ошибкой
	slog.Info("starting server on port 8080")   // НЕ должно быть ошибкой
	slog.Error("failed to connect to database") // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 2: Английский язык
	log.Println("запуск сервера")                   // ДОЛЖНО быть ошибкой
	log.Println("ошибка подключения к базе данных") // ДОЛЖНО быть ошибкой
	log.Println("starting server")                  // НЕ должно быть ошибкой
	log.Println("failed to connect to database")    // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 3: Спецсимволы и эмодзи
	log.Println("server started!🚀")                 // ДОЛЖНО быть ошибкой
	log.Println("connection failed!!!")             // ДОЛЖНО быть ошибкой
	log.Println("warning: something went wrong...") // ДОЛЖНО быть ошибкой
	log.Println("server started")                   // НЕ должно быть ошибкой
	log.Println("connection failed")                // НЕ должно быть ошибкой
	log.Println("something went wrong")             // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 4: Чувствительные данные
	password := "secret"
	apiKey := "12345"
	token := "abc"

	log.Println("user password: " + password)      // ДОЛЖНО быть ошибкой
	log.Println("api_key=" + apiKey)               // ДОЛЖНО быть ошибкой
	log.Println("token: " + token)                 // ДОЛЖНО быть ошибкой
	log.Println("user authenticated successfully") // НЕ должно быть ошибкой
	log.Println("api request completed")           // НЕ должно быть ошибкой
	log.Println("token validated")                 // НЕ должно быть ошибкой
}
