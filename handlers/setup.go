package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetupHandlers(bot *tgbotapi.BotAPI, db *db.Queries) {
	ctx := context.Background()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var chatID int64
		var userID int64
		var command string

		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
			userID = update.CallbackQuery.From.ID
			command = update.CallbackQuery.Data
		} else if update.Message != nil {
			chatID = update.Message.Chat.ID
			userID = update.Message.From.ID
			if update.Message.IsCommand() {
				command = update.Message.Command()
			}
		} else {
			continue
		}
	}
}
