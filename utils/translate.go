package utils

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

var translations = make(map[string]map[string]string)

func LoadTranslations() error {
	langs := []string{"en", "ru", "kz"}
	for _, lang := range langs {
		filePath := fmt.Sprintf("utils/languages/%s.json", lang)
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to load %s: %v", lang, err)
		}
		defer file.Close()

		var messages map[string]string
		if err := json.NewDecoder(file).Decode(&messages); err != nil {
			return fmt.Errorf("failed to parse %s: %v", lang, err)
		}

		translations[lang] = messages
	}
	return nil
}

func Translate(lang, key string) string {
	messages, langExists := translations[lang]
	if !langExists {
		log.Printf("Language '%s' not found. Falling back to English.", lang)
		messages = translations["en"]
	}

	// Check if the key exists in the language
	message, keyExists := messages[key]
	if !keyExists {
		log.Printf("Key '%s' not found in language '%s'. Falling back to English.", key, lang)
		message = translations["en"][key]
	}

	// If the key is still missing in English, return the key itself
	if message == "" {
		log.Printf("Key '%s' is missing in all languages. Returning key as message.", key)
		return key
	}

	return message
}

func GetTranslation(ctx context.Context, queries *db.Queries, update tgbotapi.Update, key string) (error, string) {
	var userID int

	if update.Message != nil {
		userID = int(update.Message.From.ID)
	} else if update.CallbackQuery != nil {
		userID = int(update.CallbackQuery.From.ID)
	}

	lang, err := queries.GetLanguage(ctx, strconv.Itoa(userID))
	if err != nil {
		log.Printf("Failed to Get Language: %v", err)
	}

	message := Translate(lang.String, key)
	return nil, message
}
