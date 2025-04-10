package handlers

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/utils"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
)

func handleCommand(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()
	username := updates.Message.From.UserName
	switch command {
	case "start":
		if updates.Message.From.UserName == "" {
			msg := tgbotapi.NewMessage(chatid, "Извините, но вы не установили в настройках Telegram юзернейм, и поэтому мы не можем вас добавить. \n\n<b>Как установить юзернейм в Telegram:</b> \n1️⃣ Откройте Telegram. \n2️⃣ Перейдите в <b>Настройки</b> (⚙️).\n3️⃣ Выберите <b>Изменить профиль</b>. \n4️⃣ Нажмите на <b>Имя пользователя</b>. \n5️⃣ Введите уникальный юзернейм.\n6️⃣ Сохраните изменения. \n После этого снова вызовите эту команду! 🚀\n\n <b>/start</b>")
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else {
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

func handleCommandAdmin(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid int64) {
	ctx := context.Background()

	switch command {
	case "write":
		msg := tgbotapi.NewMessage(chatid, "Напишите в следующем сообщении текст и он будет выслан всем пользователям!")
		bot.Send(msg)

		var state sql.NullString
		state.Valid = true
		state.String = "admin_write"
		arg := db.SetUserStateParams{State: state, Userid: strconv.FormatInt(userid, 10)}
		queries.SetUserState(ctx, arg)

	case "reg_start":
		err := queries.Enable_bot_registration(ctx)
		if err != nil {
			log.Println(err)
		}
		chatid, err = strconv.ParseInt(utils.RegisterChatID, 10, 64)
		if err != nil {
			log.Println(err)
		}
		msg := tgbotapi.NewMessage(chatid, "Registration started")
		bot.Send(msg)
	case "reg_end":
		err := queries.Diasble_bot_registration(ctx)
		if err != nil {
			log.Println(err)
		}

		msg := tgbotapi.NewMessage(chatid, "Registration ended")
		bot.Send(msg)

	case "delete_start":
		users, err := queries.GetRegisteredUsers(ctx)
		if err != nil {
			log.Println(err)
		}

		for _, user := range users {
			lng := user.Language.String
			chatID, _ := strconv.ParseInt(user.Userid, 10, 64)

			text, err := utils.GetTranslation(ctx, queries, updates, "admin_1")
			if err != nil {
				log.Println(err)
			}

			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = "HTML"

			yesNo := utils.Translate(lng, "yes_no")
			msg.ReplyMarkup = utils.InlineAccepter(strconv.FormatInt(chatID, 10), yesNo)

			sentMsg, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			} else {
				go func(chatID int64, messageID int) {
					time.AfterFunc(24*time.Hour, func() {
						deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
						if _, err := bot.Send(deleteMsg); err != nil {
							log.Println("Ошибка при удалении сообщения:", err)
						}
					})
				}(chatID, sentMsg.MessageID)
			}

			queries.DeleteUserRegedAll(ctx)
		}

		msg := tgbotapi.NewMessage(chatid, "Delete suggested")
		bot.Send(msg)

	case "delete_end":
		users, err := queries.GetUnregisteredUsers(ctx)
		if err != nil {
			log.Println(err)
		}

		for _, user := range users {
			userid, err := strconv.ParseInt(user.Userid, 10, 64)
			err = DeleteUser(queries, userid)
			if err != nil {
				log.Println(err)
			}
		}

		msg := tgbotapi.NewMessage(chatid, "Deleted all unreged")
		bot.Send(msg)
	}

}
