package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func handleCommand(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()
	username := updates.Message.From.UserName
	switch command {
	case "start":
		//todo: User created for the first time in DataBase
		arg := db.CreateUserParams{Userid: strconv.FormatInt(userid, 10), Username: username}
		err := queries.CreateUser(ctx, arg)
		if err != nil {
			log.Println(err)
		}

		//todo: ask for a language
		msg := tgbotapi.NewMessage(chatid, "Выберите язык / Тілді таңдаңыз")
		msg.ReplyMarkup = utils.InlineLang()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

		//todo: Message on start
		//key = "start_1"
		//text, err = utils.GetTranslation(ctx, queries, updates, key)
		//if err != nil {
		//	log.Println(err)
		//}
		//msg = tgbotapi.NewMessage(chatid, updates.Message.From.FirstName+text)
		//msg.ParseMode = "HTML"
		//_, err = bot.Send(msg)
		//if err != nil {
		//	log.Println(err)
		//}

	case "startorendbot":
		reged, err := queries.Getbot(ctx)
		if err != nil {
			log.Println(err)
		}
		if reged.Bool == true {
			queries.Diasble_bot_registration(ctx)
		} else {
			queries.Enable_bot_registration(ctx)
		}

	case "menu":
		key := "menu_1"
		text, err := utils.GetTranslation(ctx, queries, updates, key)
		if err != nil {
			log.Println(err)
		}
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = utils.InlineMenu()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
