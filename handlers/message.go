package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func handleMessage(message string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()
	state, err := queries.GetUserState(ctx, strconv.FormatInt(userid, 10))
	if err != nil {
		log.Println(err.Error())
	}
	if state.String == "waiting_read" {
		minutes, err := strconv.Atoi(message)
		if err != nil {
			key := "read_1"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err.Error())
			}
			msg := tgbotapi.NewMessage(chatid, text)
			msg.ParseMode = "html"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		almatyLocation, err := time.LoadLocation("Asia/Almaty")
		if err != nil {
			fmt.Println("Error loading timezone:", err)
			return
		}
		date := time.Now().In(almatyLocation)
		arg := db.CreateReadingLogParams{
			Userid:      strconv.FormatInt(userid, 10),
			Username:    updates.Message.From.UserName,
			Date:        date,
			MinutesRead: int32(minutes),
		}
		err = queries.CreateReadingLog(ctx, arg)
		if err != nil {
			log.Println("updating maybe")

			if isDuplicateKeyError(err) {
				updateArg := db.UpdateReadingLogParams{
					Userid:      strconv.FormatInt(userid, 10),
					MinutesRead: int32(minutes),
					Date:        date,
				}
				err = queries.UpdateReadingLog(ctx, updateArg)
				if err != nil {
					log.Println("UpdateReadingLog error:", err)
				}
			} else {
				key := "read_2"
				text, err := utils.GetTranslation(ctx, queries, updates, key)
				if err != nil {
					log.Println(err.Error())
				}
				msg := tgbotapi.NewMessage(chatid, text)
				msg.ParseMode = "html"
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		//todo: delete state
		err = queries.DeleteUserState(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			log.Println(err.Error())
		}
		//todo: save in sheets
		err = utils.AddReadingMinutes(utils.SheetID, strconv.FormatInt(userid, 10), minutes, date)
		if err != nil {
			log.Println(err.Error())
		}

		var key string
		if minutes > 29 {
			key = "read_3"
		} else {
			key = "read_4"
		}

		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err.Error())
		}
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "html"
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}

		time.Sleep(1 * time.Second)
		key = "menu_1"
		text, err = utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		msg = tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = utils.InlineMenu()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		//todo: say to
	} else if strings.HasPrefix(state.String, "change_read_") {
		minutes, err := strconv.Atoi(message)
		if err != nil {
			key := "read_2"
			text, err := utils.GetTranslation(ctx, queries, updates, key)
			if err != nil {
				log.Println(err.Error())
			}
			msg := tgbotapi.NewMessage(chatid, text)
			msg.ParseMode = "html"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		/*
			var text string
				_, err := fmt.Sscanf(command, "change_%v", &text)
				if err != nil {
					log.Println("Error parsing command:", err, "Command:", command)
					return
				}

				parts := strings.Split(text, ".")

				if len(parts) > 1 && parts[1] == "0" {
					parts[1] = "12"
				}

				day, _ := strconv.Atoi(parts[0])
				month, _ := strconv.Atoi(parts[1])
				year, _ := strconv.Atoi(parts[2])
				var time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

				arg := db.CreateReadingLogParams{
					Userid:      strconv.FormatInt(userid, 10),
					Username:    updates.CallbackQuery.From.UserName,
					MinutesRead: 100,
					Date:        time,
				}

				log.Println("im here", arg)
				err = queries.CreateReadingLog(ctx, arg)
				if err != nil {
					log.Println(err)
				}
		*/

		var text string
		_, err = fmt.Sscanf(state.String, "change_read_change_%v", &text)
		if err != nil {
			log.Println("Error parsing command:", err, "Command:", state.String)
			return
		}

		parts := strings.Split(text, ".")

		if len(parts) > 1 && parts[1] == "0" {
			parts[1] = "12"
		}

		day, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		year, _ := strconv.Atoi(parts[2])
		var date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		log.Println(date)

		arg := db.CreateReadingLogParams{
			Userid:      strconv.FormatInt(userid, 10),
			Username:    updates.Message.From.UserName,
			Date:        date,
			MinutesRead: int32(minutes),
		}
		err = queries.CreateReadingLog(ctx, arg)
		if err != nil {
			log.Println("updating maybe")

			if isDuplicateKeyError(err) {
				updateArg := db.UpdateReadingLogParams{
					Userid:      strconv.FormatInt(userid, 10),
					MinutesRead: int32(minutes),
					Date:        date,
				}
				err = queries.UpdateReadingLog(ctx, updateArg)
				if err != nil {
					log.Println("UpdateReadingLog error:", err)
				}
			} else {
				key := "read_2"
				text, err := utils.GetTranslation(ctx, queries, updates, key)
				if err != nil {
					log.Println(err.Error())
				}
				msg := tgbotapi.NewMessage(chatid, text)
				msg.ParseMode = "html"
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		//todo: delete state
		err = queries.DeleteUserState(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			log.Println(err.Error())
		}
		//todo: save in sheets
		err = utils.AddReadingMinutes(utils.SheetID, strconv.FormatInt(userid, 10), minutes, date)
		if err != nil {
			log.Println(err.Error())
		}
		key := "read_3"
		text, err = utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err.Error())
		}
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "html"
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
		deleteMsg := tgbotapi.NewDeleteMessage(chatid, sent.MessageID)
		time.Sleep(2 * time.Second)
		bot.Send(deleteMsg)

		userMessageID := updates.Message.MessageID
		deleteUserMsg := tgbotapi.NewDeleteMessage(chatid, userMessageID)
		_, err = bot.Request(deleteUserMsg)
		if err != nil {
			log.Println("Ошибка удаления сообщения пользователя:", err)
		}
	} else if state.String == "admin_write" {
		users, err := queries.GetRegisteredUsers(ctx)
		if err != nil {
			log.Println(err)
		}

		for _, user := range users {
			chatID, err := strconv.ParseInt(user.Userid, 10, 64)
			if err != nil {
				log.Println(err.Error())
			}
			msg := tgbotapi.NewMessage(chatID, message)
			bot.Send(msg)
		}

		queries.DeleteUserState(ctx, strconv.FormatInt(userid, 10))
	}
}

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
