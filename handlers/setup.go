package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func SetupHandlers(bot *tgbotapi.BotAPI, queries *db.Queries) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if IsBotWorking(queries, update) || (update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, "accepter")) {
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
				log.Println(chatID, userID)
				if update.Message.IsCommand() {
					command = update.Message.Command()
					log.Println(chatID, userID, command)

					if isAdmin(chatID) {
						handleCommandAdmin(command, queries, update, bot, chatID, userID)
					} else {
						handleCommand(command, queries, update, bot, chatID, userID)
					}
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

	var chatID int64
	if updates.CallbackQuery != nil {
		chatID = updates.CallbackQuery.Message.Chat.ID
	} else if updates.Message != nil {
		chatID = updates.Message.Chat.ID
	}

	reged, err := queries.Getbot(ctx)
	if err != nil {
		log.Println(err, " Setup 42 line")
		reged.Bool = false
	}

	var userReged sql.NullBool
	if strconv.FormatInt(chatID, 10) == utils.RegisterChatID {
		userReged.Bool = true
	} else {
		userReged, err = queries.GetUserReged(ctx, strconv.FormatInt(chatID, 10))
		if err != nil {
			log.Println(err, " Setup 55 line")
			userReged.Bool = false
		}
	}

	return reged.Bool || userReged.Bool
}

func SendReminders(bot *tgbotapi.BotAPI, queries *db.Queries) {
	ctx := context.Background()
	users, err := queries.GetUsersWithoutReadingToday(ctx)
	if err != nil {
		log.Println(err)
	}

	for _, user := range users {
		chatID, err := strconv.Atoi(user.Userid)
		if err != nil {
			log.Println(err)
		}
		text := utils.Translate(user.Language.String, "remind_1")
		msg := tgbotapi.NewMessage(int64(chatID), text)
		msg.ParseMode = "HTML"
		bot.Send(msg)
	}
}

func isAdmin(chatID int64) bool {
	AdminID, err := strconv.ParseInt(utils.RegisterChatID, 10, 64)
	if err != nil {
		log.Println(err)
	}
	if AdminID == chatID {
		return true
	}
	return false
}
