package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	apiURL := "https://petstore.swagger.io/v2/pet/findByStatus?status=available"

	// Выполняем HTTP-запрос GET к API
	response, err := http.Get(apiURL)
	if err != nil {
		log.Fatalf("Не удалось выполнить HTTP-запрос: %v", err)
	}
	defer response.Body.Close()

	// Проверяем статус ответа
	if response.StatusCode != http.StatusOK {
		log.Fatalf("HTTP-запрос завершился неудачно: %s", response.Status)
	}

	// Декодируем JSON-ответ в структуру данных
	var pets []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(response.Body).Decode(&pets); err != nil {
		log.Fatalf("Не удалось декодировать JSON: %v", err)
	}

	// Выводим информацию о питомцах
	fmt.Println("Доступные питомцы:", pets)
}
