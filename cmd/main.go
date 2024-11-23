package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

const telegramToken = "8057819715:AAHbLvhxrvhldtyQAS2y7vFzqH_GflyO0nI"

func generateImage(prompt string) (string, error) {
	url := "http://localhost:7860/sdapi/v1/txt2img"

	// Формируем данные для запроса
	requestBody := map[string]interface{}{
		"negative_prompt":     "Low quality, nude",
		"sampler_name":        "Euler a",
		"prompt":              prompt,
		"width":               512,
		"height":              768,
		"num_inference_steps": 50,
		"cfg_scale":           12,
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

	// Логируем полученный ответ для диагностики
	log.Printf("Ответ от сервера: %v", response)

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

func main() {
	// Настроим логирование
	logFile, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}
	defer logFile.Close()

	// Создаем логгер, который будет записывать в файл и выводить в консоль
	multiLogger := log.New(os.Stdout, "", log.LstdFlags)
	log.SetOutput(logFile)           // Устанавливаем вывод в файл
	multiLogger.SetOutput(os.Stdout) // Устанавливаем вывод в консоль

	// Создаем нового бота
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		multiLogger.Fatalf("Не удалось авторизовать бота: %v", err)
	}
	multiLogger.Printf("Авторизован как %s", bot.Self.UserName)

	// Получаем обновления от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		multiLogger.Fatalf("Не удалось получить обновления от бота: %v", err)
	}

	// Создаем клавиатуру с кнопкой "Сгенерировать машину"
	button := tgbotapi.NewInlineKeyboardButtonData("Сгенерировать девушку", "generate_car")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)

	for update := range updates {
		// Обрабатываем обычные сообщения
		if update.Message != nil && update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				multiLogger.Printf("Пользователь %s начал чат с ботом", update.Message.From.UserName)

				// Отправляем кнопку с командой "Сгенерировать машину"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите кнопку для генерации машины!")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}
		}

		// Обрабатываем нажатие на кнопку
		if update.CallbackQuery != nil {
			multiLogger.Printf("Пользователь %s нажал на кнопку: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

			if update.CallbackQuery.Data == "generate_car" {
				// Генерируем изображение
				multiLogger.Println("Запрос на генерацию изображения: красная машина")
				imgData, err := generateImage("beautiful woman") // Пример промпта
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при генерации изображения"))
					multiLogger.Printf("Ошибка при генерации изображения: %v", err)
					continue
				}

				// Логируем base64 данные для диагностики
				multiLogger.Printf("Полученные base64 данные изображения: %s", imgData)

				// Декодируем base64 строку изображения
				decodedImage, err := base64.StdEncoding.DecodeString(imgData)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при декодировании изображения"))
					multiLogger.Printf("Ошибка при декодировании изображения: %v", err)
					continue
				}

				// Сохраняем изображение во временный файл
				tmpFile, err := os.CreateTemp("", "generated_image_*.png")
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при создании файла"))
					multiLogger.Printf("Ошибка при создании файла: %v", err)
					continue
				}
				defer tmpFile.Close()

				if _, err := tmpFile.Write(decodedImage); err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при записи изображения"))
					multiLogger.Printf("Ошибка при записи изображения: %v", err)
					continue
				}

				// Отправляем изображение как файл
				photo := tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, tmpFile.Name())
				bot.Send(photo)

				// Ответ на нажатие кнопки
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Генерация завершена!")
				bot.AnswerCallbackQuery(callback)
				multiLogger.Println("Изображение отправлено пользователю.")
			}
		}
	}
}
