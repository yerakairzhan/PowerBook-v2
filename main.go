package main

import (
	db "PowerBook2.0/db/sqlc"
	"PowerBook2.0/handlers"
	"PowerBook2.0/utils"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	utils.LoadConfig()
	err := utils.LoadTranslations()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	bot, err := tgbotapi.NewBotAPI(utils.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	db := db.New(conn)
	db.CreateBot(context.Background())

	// Инициализируем планировщик задач (cron)
	c := cron.New()

	// Добавляем ежедневное напоминание в 20:00
	c.AddFunc("0 20 * * *", func() {
		handlers.SendReminders(bot, db)
	})

	c.Start()

	handlers.SetupHandlers(bot, db)

	select {}
}
