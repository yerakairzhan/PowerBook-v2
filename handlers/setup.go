package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func SetupHandlers(bot *tgbotapi.BotAPI, queries *db.Queries) {
	//ctx := context.Background()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if IsBotWorking(queries, update) {
			var userID int64
			var chatID int64
			var command string

			if update.CallbackQuery != nil {
				chatID = update.CallbackQuery.Message.Chat.ID
				userID = update.CallbackQuery.From.ID
				command = update.CallbackQuery.Data

				handleCallback(command, queries, update, bot, chatID, userID)

			} else if update.Message != nil {
				chatID = update.Message.Chat.ID
				userID = update.Message.From.ID
				if update.Message.IsCommand() {
					command = update.Message.Command()
					handleCommand(command, queries, update, bot, chatID, userID)
				} else {
					message := update.Message.Text
					handleMessage(message, queries, update, bot, chatID, userID)
				}
			} else {
				continue
			}
		} else {
			text := utils.Translate("ru", "bot_1")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		}
	}
}

func IsBotWorking(queries *db.Queries, updates tgbotapi.Update) bool {
	utils.LoadConfig()
	ctx := context.Background()

	var userid int64
	if updates.CallbackQuery != nil {
		userid = updates.CallbackQuery.From.ID
	} else if updates.Message != nil {
		userid = updates.Message.From.ID
	}

	reged, err := queries.Getbot(ctx)
	if err != nil {
		log.Println(err, " Setup 42 line")
		reged.Bool = false
	}

	var userReged sql.NullBool
	if strconv.FormatInt(userid, 10) == utils.AdminID {
		userReged.Bool = true
	} else {
		userReged, err = queries.GetUserReged(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			log.Println(err, " Setup 55 line")
			userReged.Bool = false
		}
	}
	return reged.Bool || userReged.Bool
}
