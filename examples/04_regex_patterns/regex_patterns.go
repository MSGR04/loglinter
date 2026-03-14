package main

import "log"

func main() {
	log.Println("confidential report")  // должно быть ошибкой
	log.Println("secret_project_alpha") // должно быть ошибкой (secret_.*)
	log.Println("secret_data_123")      // должно быть ошибкой (secret_.*)
	log.Println("internal_456")         // должно быть ошибкой (internal_[0-9]+)
	log.Println("internal_abc")         // НЕ должно быть ошибкой (не подходит под паттерн)
}
