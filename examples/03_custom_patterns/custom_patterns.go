package main

import "log"

func main() {
	// ДОЛЖНЫ быть ошибками (кастомные паттерны)
	log.Println("This is confidential data")
	log.Println("internal_only: 12345")
	log.Println("secret_project: alpha")

	// НЕ ДОЛЖНЫ быть ошибками (стандартные отключены)
	log.Println("password: 123")
	log.Println("api_key: abc")
	log.Println("token: xyz")

	// Другие правила продолжают работать
	log.Println("Starting server") // ошибка о заглавной
}
