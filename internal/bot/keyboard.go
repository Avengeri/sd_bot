package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// GenerateKeyboard создает клавиатуру с кнопкой "Сгенерировать изображение"
func GenerateKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Сгенерировать изображение", "generate_image")
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)
}
