package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"sd_bot/internal/bot"
)

func main() {

	// Загрузка переменных окружения из .env файла в корне проекта
	err := godotenv.Load("./go.env")
	if err != nil {
		log.Println("Файл .env не найден или не удалось его загрузить. Переменные окружения будут загружены из системы.")
	}

	// Получаем значения переменных окружения
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	imageAPIURL := os.Getenv("IMAGE_API_URL")

	if telegramToken == "" || imageAPIURL == "" {
		log.Fatal("Отсутствуют необходимые переменные окружения (TELEGRAM_TOKEN, IMAGE_API_URL).")
	}

	// Запуск Telegram-бота
	bot.Run(telegramToken, imageAPIURL)
}
