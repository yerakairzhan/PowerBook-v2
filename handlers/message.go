package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
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

		date := time.Now()
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

		//todo: say to
	}
}

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
