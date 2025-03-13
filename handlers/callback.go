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
	"time"
)

func handleCallback(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()
	log.Println("callback: ", command)
	switch {
	case strings.HasPrefix(command, "lang"):
		var lang sql.NullString
		lang.String = strings.TrimPrefix(command, "lang_")
		lang.Valid = true
		arg := db.SetLanguageParams{Language: lang, Userid: strconv.FormatInt(userid, 10)}
		err := queries.SetLanguage(ctx, arg)
		if err != nil {
			log.Println(err)
		}

		//todo: Change the text for Start
		key := "start_1"
		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		inlineKeyboard := utils.InlineRegister()
		callback := updates.CallbackQuery

		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			callback.From.FirstName+text,
			inlineKeyboard,
		)

		_, err = bot.Send(editMsg)
		if err != nil {
			log.Println("Ошибка при изменении сообщения:", err)
		}

	case command == "register":
		//todo: Change the text for waiting for accepts sending To admin

		chatidAdmin, err := strconv.ParseInt(utils.RegisterChatID, 10, 64)
		if err != nil {
			log.Println(err)
		}
		now := time.Now()
		text := "Who : @" + updates.CallbackQuery.From.UserName + "\nWhen : " + now.Local().Format("2006-01-02 15:04:05") + "\n(make sure choosing below!)"
		msg := tgbotapi.NewMessage(chatidAdmin, text)
		inlineKeyboard := utils.InlineAccepter(strconv.FormatInt(chatid, 10))
		msg.ReplyMarkup = inlineKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		messageID := updates.CallbackQuery.Message.MessageID
		removeInlineButtons(bot, chatid, messageID)

	case strings.HasPrefix(command, "accepter"):
		trimmed := strings.TrimPrefix(command, "accepter_")
		parts := strings.Split(trimmed, "_")
		choice := parts[0]
		chatID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
		}

		if choice == "yes" {
			key := "register_1"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err)
			}
			msg := tgbotapi.NewMessage(chatID, text)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

			//todo: Save in db
			err = queries.SetUserReged(ctx, strconv.FormatInt(userid, 10))
			if err != nil {
				log.Println(err)
			}

			//todo: Send the instructions of the bot
			key = "start_2"
			text, err = utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err)
			}
			msg = tgbotapi.NewMessage(chatID, text)

			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else if choice == "no" {
			key := "register_2"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err)
			}
			msg := tgbotapi.NewMessage(chatID, text)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
		messageID := updates.CallbackQuery.Message.MessageID
		removeInlineButtons(bot, chatid, messageID)

	case strings.HasPrefix(command, "menu"):
		command = strings.TrimPrefix(command, "menu_")
		if command == "read" {
			//todo: ask for minutes
			key := "read_1"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err)
			}
			msg := tgbotapi.NewMessage(chatid, text)
			msg.ParseMode = "HTML"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			var state sql.NullString
			state.Valid = true
			state.String = "waiting_read"
			arg := db.SetUserStateParams{
				Userid: strconv.FormatInt(userid, 10),
				State:  state,
			}
			err = queries.SetUserState(ctx, arg)
			if err != nil {
				log.Println(err)
			}
		}
	case strings.HasPrefix(command, "calendar"):
		trimmed := strings.TrimPrefix(command, "calendar_")
		parts := strings.Split(trimmed, "_")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		messageID := updates.CallbackQuery.Message.MessageID

		sendCalendar(year, month, queries, updates, bot, chatid, userid, true, messageID)
	}
}

func removeInlineButtons(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	removeKeyboard := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	})

	_, err := bot.Send(removeKeyboard)
	if err != nil {
		log.Println("Ошибка при удалении кнопок:", err)
	}
}

func sendCalendar(year int, month int, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64, isEdit bool, messageID int) {
	ctx := context.Background()

	readLogs, err := queries.GetReadingLogsByUser(ctx, strconv.FormatInt(userid, 10))
	if err != nil {
		log.Println(err.Error())
	}

	readMinutes := make(map[int]int)

	for _, log := range readLogs {
		if int(log.Date.Month()) == int(month) && log.Date.Year() == year {
			day := log.Date.Day()
			readMinutes[day] = int(log.MinutesRead)
		}
	}

	key := "experience_1"
	text, err := utils.GetTranslation(ctx, queries, updates, key)
	if err != nil {
		log.Println(err)
	}

	inline := utils.InlineCalendarKeyboard(year, int(month), readMinutes)

	if isEdit {
		editMsg := tgbotapi.NewEditMessageText(chatid, messageID, text)
		editMsg.ReplyMarkup = &inline
		editMsg.ParseMode = "HTML"
		bot.Send(editMsg)
	} else {
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ReplyMarkup = inline
		msg.ParseMode = "HTML"
		bot.Send(msg)
	}
}
