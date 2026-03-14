package main

import "log"

func main() {
	// ДОЛЖНО быть ошибкой (русский язык)
	log.Println("запуск сервера")

	// НЕ ДОЛЖНО быть ошибкой (lowercase отключен в конфиге)
	log.Println("Starting server")

	// ДОЛЖНО быть ошибкой (чувствительные данные - есть в паттернах)
	log.Println("my_secret_value: 123")

	// ДОЛЖНО быть ошибкой (спецсимволы)
	log.Println("server started!!!")

	// НЕ ДОЛЖНО быть ошибкой (нет в паттернах, стандартные отключены)
	log.Println("password: 12345")

	// НЕ ДОЛЖНО быть ошибкой (нет в паттернах)
	log.Println("api_key: 67890")
}
