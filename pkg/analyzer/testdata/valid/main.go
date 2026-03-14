package main

import (
	"log"
	"log/slog"
)

func main() {
	// Правильные примеры
	log.Println("starting server")
	log.Printf("server running on port %d", 8080)
	slog.Info("user authenticated")

	// Неправильные примеры (для тестирования)
	log.Println("Starting server")         // должно быть с маленькой буквы
	slog.Info("запуск сервера")            // не английский
	log.Println("server started!🚀")        // спецсимволы и эмодзи
	log.Printf("user password: %s", "123") // чувствительные данные
}
