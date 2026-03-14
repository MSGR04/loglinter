package _2_config

import (
	"log"
	"log/slog"
)

func main() {
	// Эти сообщения должны проверяться (все правила включены)
	log.Println("Starting server")       // ДОЛЖНО быть ошибкой (заглавная)
	slog.Info("запуск сервера")          // ДОЛЖНО быть ошибкой (русский)
	log.Println("server started!🚀")      // ДОЛЖНО быть ошибкой (спецсимволы)
	log.Println("user password: secret") // ДОЛЖНО быть ошибкой (чувствительные данные)

	// Эти сообщения правильные
	log.Println("starting server") // НЕ должно быть ошибкой
	slog.Info("server started")    // НЕ должно быть ошибкой
}
