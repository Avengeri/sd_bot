package imggen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Функция для генерации изображения через API
func GenerateImage(prompt string, imageAPIURL string) (string, error) {
	url := imageAPIURL // URL API для генерации изображений

	// Формируем данные для запроса
	requestBody := map[string]interface{}{
		"prompt":              prompt,
		"width":               512,
		"height":              512,
		"num_inference_steps": 30,
	}

	// Преобразуем запрос в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Ошибка при маршалинге данных: %v", err)
		return "", err
	}

	// Отправляем POST запрос на сервер
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Ошибка при отправке запроса: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка ответа от сервера: %v", resp.Status)
		return "", fmt.Errorf("Ошибка ответа от сервера: %v", resp.Status)
	}

	// Читаем и декодируем ответ
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Ошибка при декодировании ответа: %v", err)
		return "", err
	}

	// Проверяем наличие данных "images"
	images, ok := response["images"].([]interface{})
	if !ok {
		log.Printf("Не удалось извлечь изображения из ответа: %v", response["images"])
		return "", fmt.Errorf("Не удалось извлечь изображения из ответа")
	}

	// Если изображение найдено, возвращаем base64 строку изображения
	if len(images) > 0 {
		imageData, ok := images[0].(string)
		if !ok {
			log.Printf("Ошибка при преобразовании изображения: %v", images[0])
			return "", fmt.Errorf("Ошибка при преобразовании изображения")
		}
		return imageData, nil
	}

	// Если изображений нет, возвращаем ошибку
	log.Println("Изображение не было сгенерировано.")
	return "", fmt.Errorf("Изображение не было сгенерировано")
}
