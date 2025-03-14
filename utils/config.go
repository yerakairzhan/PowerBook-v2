package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var BotToken string
var DBDriver string
var DBSource string
var AdminID string
var RegisterChatID string
var TableURL string
var SheetID string
var Creds string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла : ", err)
	}

	BotToken = os.Getenv("BOT_TOKEN")
	DBDriver = os.Getenv("DB_DRIVER")
	DBSource = os.Getenv("DB_SOURCE")
	TableURL = os.Getenv("TABLE_URL")
	SheetID = os.Getenv("GOOGLE_API")
	Creds = os.Getenv("GOOGLE_CREDENTIALS")
	AdminID = os.Getenv("ADMIN_ID")
	RegisterChatID = os.Getenv("REGISTER_CHAT_ID")
}
