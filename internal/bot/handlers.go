package bot

import (
	"encoding/base64"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"sd_bot/internal/imggen"
)

// Функция для обработки нажатия на кнопку
func HandleCallbackQuery(update tgbotapi.Update, botAPI *tgbotapi.BotAPI, imageAPIURL string) {
	// Логируем нажатие на кнопку
	log.Printf("Пользователь %s нажал на кнопку: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

	if update.CallbackQuery.Data == "generate_image" {
		// Генерируем изображение
		log.Println("Запрос на генерацию изображения: описание для изображения")
		imageData, err := imggen.GenerateImage("Описание для генерации изображения", imageAPIURL) // Например, передаем описание
		if err != nil {
			botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при генерации изображения"))
			log.Printf("Ошибка при генерации изображения: %v", err)
			return
		}

		// Логируем успешную генерацию
		log.Println("Изображение сгенерировано успешно")

		// Преобразуем строку base64 в бинарные данные
		imageBytes, err := decodeBase64Image(imageData)
		if err != nil {
			botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при обработке изображения"))
			log.Printf("Ошибка при обработке base64 изображения: %v", err)
			return
		}

		// Создаем файл с изображением во временном файле
		tempFile, err := ioutil.TempFile("", "image_*.jpg")
		if err != nil {
			botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при создании файла изображения"))
			log.Printf("Ошибка при создании временного файла: %v", err)
			return
		}
		defer tempFile.Close()

		// Записываем бинарные данные изображения в файл
		_, err = tempFile.Write(imageBytes)
		if err != nil {
			botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при записи файла изображения"))
			log.Printf("Ошибка при записи изображения в файл: %v", err)
			return
		}

		// Отправляем изображение как файл
		photo := tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, tempFile.Name())
		_, err = botAPI.Send(photo)
		if err != nil {
			log.Printf("Ошибка при отправке изображения: %v", err)
			botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при отправке изображения"))
		}

		// Отправляем ответ на callback
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Генерация завершена!")
		botAPI.AnswerCallbackQuery(callback)
	}
}

// Функция для декодирования base64 строки в бинарные данные
func decodeBase64Image(base64String string) ([]byte, error) {
	// Убираем префикс "data:image/jpeg;base64," если он есть
	if len(base64String) > 22 && base64String[:22] == "data:image/jpeg;base64," {
		base64String = base64String[22:]
	}

	// Декодируем строку Base64 в бинарные данные
	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования base64: %v", err)
	}

	return imageBytes, nil
}
