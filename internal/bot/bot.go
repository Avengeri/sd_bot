package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Run(telegramToken, imageAPIURL string) {
	// Создаем нового бота
	botAPI, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("Не удалось авторизовать бота: %v", err)
	}

	log.Printf("Авторизован как %s", botAPI.Self.UserName)

	// Получаем обновления от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Не удалось получить обновления от бота: %v", err)
	}

	// Обработка обновлений
	for update := range updates {
		// Обрабатываем обычные сообщения
		if update.Message != nil {
			if update.Message.IsCommand() {
				if update.Message.Command() == "start" {
					log.Printf("Пользователь %s начал чат с ботом", update.Message.From.UserName)

					// Отправляем клавиатуру с кнопкой "Сгенерировать изображение"
					keyboard := GenerateKeyboard()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите кнопку для генерации изображения!")
					msg.ReplyMarkup = keyboard
					botAPI.Send(msg)
				}
			}
		}

		// Обрабатываем нажатие на кнопку
		if update.CallbackQuery != nil {
			HandleCallbackQuery(update, botAPI, imageAPIURL)
		}
	}
}
