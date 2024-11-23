package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sd_bot/internal/constans"
	"sd_bot/internal/model"
)

func newKeyboardRow(buttonTexts ...string) []tgbotapi.KeyboardButton {
	var buttons []tgbotapi.KeyboardButton

	for _, text := range buttonTexts {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(text))
	}
	return buttons
}

func newInlineKeyboard(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func ShowStartMessage(bot *tgbotapi.BotAPI, user *model.User) {

	msg := tgbotapi.NewMessage(user.ChatId, "Есть аккаунт-авторизуйся, нет - заведи")
	replyRow := newKeyboardRow(constans.BUTTON_REPLY_TEXT_AUTHORIZE)

	replyKeyboard := tgbotapi.NewReplyKeyboard(replyRow)

	msg.ReplyMarkup = replyKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Не удалось отправить сообщение")
	}
}
