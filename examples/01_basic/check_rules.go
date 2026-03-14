package main

import (
	"log"
	"log/slog"
)

func main() {
	// ТЕСТ ПРАВИЛА 1: Строчная буква
	log.Info("Starting server on port 8080")    // ДОЛЖНО быть ошибкой
	slog.Error("Failed to connect to database") // ДОЛЖНО быть ошибкой
	log.Info("starting server on port 8080")    // НЕ должно быть ошибкой
	slog.Error("failed to connect to database") // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 2: Английский язык
	log.Info("запуск сервера")                    // ДОЛЖНО быть ошибкой
	log.Error("ошибка подключения к базе данных") // ДОЛЖНО быть ошибкой
	log.Info("starting server")                   // НЕ должно быть ошибкой
	log.Error("failed to connect to database")    // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 3: Спецсимволы и эмодзи
	log.Info("server started!🚀")                 // ДОЛЖНО быть ошибкой
	log.Error("connection failed!!!")            // ДОЛЖНО быть ошибкой
	log.Warn("warning: something went wrong...") // ДОЛЖНО быть ошибкой
	log.Info("server started")                   // НЕ должно быть ошибкой
	log.Error("connection failed")               // НЕ должно быть ошибкой
	log.Warn("something went wrong")             // НЕ должно быть ошибкой

	// ТЕСТ ПРАВИЛА 4: Чувствительные данные
	password := "secret"
	apiKey := "12345"
	token := "abc"

	log.Info("user password: " + password)      // ДОЛЖНО быть ошибкой
	log.Debug("api_key=" + apiKey)              // ДОЛЖНО быть ошибкой
	log.Info("token: " + token)                 // ДОЛЖНО быть ошибкой
	log.Info("user authenticated successfully") // НЕ должно быть ошибкой
	log.Debug("api request completed")          // НЕ должно быть ошибкой
	log.Info("token validated")                 // НЕ должно быть ошибкой
}
