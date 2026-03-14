package mylogger

import (
	"log"
	"log/slog"
)

func TestLogMessages() {
	// Эти сообщения должны проходить проверку
	log.Println("starting server")        // OK
	slog.Info("user authenticated")       // OK
	log.Printf("processed %d items", 100) // OK

	// Эти сообщения должны вызывать ошибки
	log.Println("Starting server")      // want `лог-сообщение должно начинаться со строчной буквы`
	slog.Info("запуск сервера")         // want `лог-сообщение должно содержать только английские символы`
	log.Println("server started!🚀")     // want `лог-сообщение не должно содержать спецсимволы или эмодзи`
	log.Println("user password: 12345") // want `лог-сообщение не должно содержать потенциально чувствительные данные`
}
