package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
)

func main() {

	//if err := godotenv.Load("./.env"); err != nil {
	//	log.Fatalf("error loading environment variable: %s", err.Error())
	//}

	//db, err := postgres.NewPostgresDB(postgres.Config{
	//	Host:     os.Getenv("DB_HOST"),
	//	Port:     os.Getenv("DB_PORT"),
	//	Username: os.Getenv("DB_USER"),
	//	DBName:   os.Getenv("DB_NAME"),
	//	SSLMode:  os.Getenv("DB_SSLMODE"),
	//	Password: os.Getenv("DB_PASSWORD"),
	//})
	//if err != nil {
	//	log.Fatalf("Postgres creation error: %s", err.Error())
	//}
	//ok, err := postgres.CheckDBConn(db)
	//if err != nil {
	//	log.Fatalf("Connection to the database could not be established: %s", err.Error())
	//}
	//fmt.Println(ok)

	//repo := repository.NewStorageUserPostgres(db)
	//service := service2.NewUserService(repo)
	//fmt.Println(service)

	bot, err := tgbotapi.NewBotAPI("8057819715:AAHbLvhxrvhldtyQAS2y7vFzqH_GflyO0nI")
	if err != nil {
		log.Fatalf("Connection to BOT API failed: %s", err.Error())
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Could not get updates from chan: %s", err.Error())
	}

	for update := range updates {
		if update.Message != nil { // Если получено сообщение
			handleMessage(bot, &update)
		} else if update.CallbackQuery != nil { // Если нажата кнопка
			handleButtonPress(bot, update.CallbackQuery)
		}
	}
}

// Обработка текстовых сообщений
func handleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Отправляем сообщение с кнопкой
	button := tgbotapi.NewInlineKeyboardButtonData("Сгенерировать машинку", "send_request")
	row := tgbotapi.NewInlineKeyboardRow(button)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сгенерируем картинку?")
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send message: %s", err.Error())
	}
}

// Обработка нажатия кнопки
func handleButtonPress(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	if callback.Data == "send_request" {
		// Отправляем запрос на внешний сервис
		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1") // Замените на нужный URL
		if err != nil {
			log.Printf("Failed to send request: %s", err.Error())
			_, _ = bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка отправки запроса."))
			return
		}
		defer resp.Body.Close()

		// Ответ пользователю
		reply := tgbotapi.NewMessage(callback.Message.Chat.ID, "Запрос отправлен успешно!")
		_, _ = bot.Send(reply)
	}

	// Уведомление Telegram о завершении обработки кнопки
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Запрос обрабатывается..."))
}
