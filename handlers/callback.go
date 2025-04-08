package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

var PhotoStartKaz string = "https://i.imghippo.com/files/fit6422iZY.png"
var PhotoStartRus string = "https://i.imghippo.com/files/ur5705pvI.png"

func handleCallback(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()
	log.Println("callback: ", command)
	switch {
	case command == "back":
		key := "menu_1"
		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		callback := updates.CallbackQuery
		inlineKeyboard := utils.InlineMenu()
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			text,
			inlineKeyboard,
		)
		editMsg.ParseMode = "HTML"
		_, err = bot.Send(editMsg)
		if err != nil {
			log.Println("Ошибка при изменении сообщения:", err)
		}

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
		// Determine the correct photo URL
		photoURL := PhotoStartKaz // Default to Kazakh
		if lang.String == "Ru" {
			photoURL = PhotoStartRus
		}

		// Create a photo message
		photo := tgbotapi.NewPhoto(chatid, tgbotapi.FileURL(photoURL))

		// Get translated text
		key := "start_1"
		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}

		// Create inline keyboard
		inlineKeyboard := utils.InlineRegister()

		// Set caption and attach inline keyboard
		photo.Caption = updates.CallbackQuery.From.FirstName + text
		photo.ReplyMarkup = inlineKeyboard
		photo.ParseMode = "HTML"

		// Send the photo with caption and inline keyboard
		_, err = bot.Send(photo)
		if err != nil {
			log.Println("Ошибка при отправке фото:", err)
		}

	case command == "register":
		//todo: Change the text for waiting for accepts sending To admin
		messageID := updates.CallbackQuery.Message.MessageID
		removeInlineButtons(bot, chatid, messageID)

		key := "register_1"
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

		//todo: Save in db
		err = queries.SetUserReged(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			log.Println(err)
		}
		//sheets
		if err := utils.AddUserToSheet(utils.SheetID, strconv.FormatInt(userid, 10), updates.CallbackQuery.From.UserName); err != nil {
			log.Fatalf("Error adding user to sheet: %v", err)
		}

		//todo: Send the instructions of the bot
		key = "start_2"
		text, err = utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		msg = tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"

		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

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
			msg.ParseMode = "HTML"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

		} else if choice == "no" {
			err := DeleteUser(queries, updates, bot, chatID, userid)
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
		} else if command == "standings" {
			leaderboard, err := queries.GetReadingLeaderboard(ctx)

			if err != nil {
				log.Println(err)
			}
			YourMax, err := queries.GetSumReading(ctx, strconv.FormatInt(userid, 10))
			if err != nil {
				log.Fatal("Error getting leaderboard:", err)
			}
			inline := utils.InlineLeaderboard(leaderboard, YourMax)
			key := "standings_1"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err)
			}
			callback := updates.CallbackQuery.Message.MessageID
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				chatid,
				callback,
				text,
				inline,
			)
			msg.ParseMode = "HTML"
			_, err = bot.Send(msg)
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

		sendCalendar(true, year, month, queries, updates, bot, chatid, userid, true, messageID)

	case strings.HasPrefix(command, "back_"):
		trimmed := strings.TrimPrefix(command, "back_")
		parts := strings.Split(trimmed, "_")
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		messageID := updates.CallbackQuery.Message.MessageID

		sendCalendar(false, year, month, queries, updates, bot, chatid, userid, true, messageID)

	case strings.HasPrefix(command, "day"):
		var text string
		_, err := fmt.Sscanf(command, "day_%v", &text)
		if err != nil {
			log.Println("Error parsing command:", err, "Command:", command)
			return
		}

		parts := strings.Split(text, ".")

		if len(parts) > 1 && parts[1] == "0" {
			parts[1] = "12"
		}

		output := strings.Join(parts, ".")
		callbackQueryID := updates.CallbackQuery.ID
		callback := tgbotapi.NewCallback(callbackQueryID, output)
		bot.Request(callback)

	case strings.HasPrefix(command, "change"): // change_17.4.2025
		var state sql.NullString
		state.Valid = true
		state.String = "change_read_" + command
		arg := db.SetUserStateParams{
			Userid: strconv.FormatInt(userid, 10),
			State:  state,
		}
		err := queries.SetUserState(ctx, arg)
		if err != nil {
			log.Println(err)
		}
		key := "change_1"
		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

		deleteMsg := tgbotapi.NewDeleteMessage(chatid, sent.MessageID)
		time.Sleep(3 * time.Second)
		bot.Request(deleteMsg)
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

func sendCalendar(simple_calendar bool, year int, month int, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64, isEdit bool, messageID int) {
	if month < 0 {
		month += 12
		year--
	} else if month >= 12 {
		month -= 12
		year++
	}

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
	var inline tgbotapi.InlineKeyboardMarkup
	var key string
	if simple_calendar {
		inline = utils.InlineCalendarKeyboard(year, int(month), readMinutes)
		key = "experience_1"
	} else {
		inline = utils.InlineCalendarChanger(year, int(month), readMinutes)
		key = "changer_1"
	}
	text, err := utils.GetTranslation(ctx, queries, updates, key)
	if err != nil {
		log.Println(err)
	}

	if isEdit {
		editMsg := tgbotapi.NewEditMessageText(chatid, messageID, text)
		editMsg.ReplyMarkup = &inline
		editMsg.ParseMode = "HTML"
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Println(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ReplyMarkup = inline
		msg.ParseMode = "HTML"
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func DeleteUser(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) error {
	ctx := context.Background()

	key := "register_2"
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

	if err := utils.DeleteUserFromSheet(utils.SheetID, strconv.FormatInt(userid, 10)); err != nil {
		log.Fatalf("Error deleting user to sheet: %v", err)
	}

	err = queries.DeleteUserReged(ctx, strconv.FormatInt(userid, 10))
	if err != nil {
		log.Println(err)
	}

	return err
}
