package utils

import (
	db "PowerBook2.0/db/sqlc"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
)

func InlineLang() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇰🇿 Қазақша	", "lang_kz"),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "callback_ru"),
		//),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "callback_en"),
		//),
	)
	return inline
}

func InlineRegister() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅", "register"),
		),
	)
	return inline
}

func InlineAccepter(chatid string, yes_no string) tgbotapi.InlineKeyboardMarkup {
	yes := yes_no[:3]
	no := yes_no[5:]
	log.Println(yes, no)

	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(yes, "accepter_yes_"+chatid),
			tgbotapi.NewInlineKeyboardButtonData(no, "accepter_no_"+chatid),
		),
	)
	return inline
}

func InlineMenu() tgbotapi.InlineKeyboardMarkup {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Read", "menu_read"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📊 Table", TableURL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💪 Progress", fmt.Sprintf("calendar_%d_%d", year, month)),
			tgbotapi.NewInlineKeyboardButtonData("🏆 Top", "menu_standings"),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("🏆 Top", "menu_standings"),
		//),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏪ Archive", fmt.Sprintf("back_%d_%d", year, month)),
		),
	)
	return inline
}

func InlineCalendarKeyboard(year int, month int, readMinutes map[int]int) tgbotapi.InlineKeyboardMarkup {
	daysOfWeek := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	months := []string{"December", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November"}
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var weekRow []tgbotapi.InlineKeyboardButton

	for _, day := range daysOfWeek {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}
	keyboard = append(keyboard, weekRow)

	// Get first weekday and total days in month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	startWeekday := int(firstDay.Weekday()) // Sunday = 0
	if startWeekday == 0 {
		startWeekday = 7 // Adjust for Monday start
	}
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	var row []tgbotapi.InlineKeyboardButton
	// Fill empty slots before the first day
	for i := 1; i < startWeekday; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}

	// Fill in the actual days
	for day := 1; day <= daysInMonth; day++ {
		minutes := readMinutes[day]

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(minutes), fmt.Sprintf("day_%d.%d.%d", day, month, year)))

		// Break at the end of each week
		if len(row) == 7 {
			keyboard = append(keyboard, row)
			row = nil
		}
	}

	// Add remaining row if not complete
	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
		}
		keyboard = append(keyboard, row)
	}

	navRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("calendar_%d_%d", year, month-1)),
		tgbotapi.NewInlineKeyboardButtonData("📆"+months[month], fmt.Sprintf("calendar_%d_%d", time.Now().Year(), int(time.Now().Month()))),
		tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("calendar_%d_%d", year, month+1)),
	}
	keyboard = append(keyboard, navRow)
	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙️", "back")))

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func InlineLeaderboard(items []db.GetReadingLeaderboardRow, usersMax db.GetSumReadingRow) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	inTop := false
	log.Println(items, usersMax)

	navRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Place", "ignore"),
		tgbotapi.NewInlineKeyboardButtonData("Username", "ignore"),
		tgbotapi.NewInlineKeyboardButtonData("Points", "ignore"),
		tgbotapi.NewInlineKeyboardButtonData("Minutes", "ignore"),
	}
	keyboard = append(keyboard, navRow)

	for i := 0; i < len(items); i++ {
		if items[i].Userid == usersMax.Userid {
			inTop = true
		}
		var emoji string
		if i == 0 {
			emoji = "🥇"
		} else if i == 1 {
			emoji = "🥈"
		} else if i == 2 {
			emoji = "🥉"
		} else if i == 3 {
			emoji = "4️⃣"
		} else if i == 4 {
			emoji = "5️⃣"
		}

		navRow := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(emoji, "ignore"),
			tgbotapi.NewInlineKeyboardButtonData("@"+items[i].Username, "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d балл.", items[i].DaysReadMoreThan30), "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d мин.", items[i].TotalMinutes), "ignore"),
		}
		keyboard = append(keyboard, navRow)
	}
	if !inTop {
		navRow := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🫵", "ignore"),
			tgbotapi.NewInlineKeyboardButtonData("@"+usersMax.Username, "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d балл.", usersMax.DaysReadMoreThan30), "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d мин.", usersMax.Sum), "ignore"),
		}
		keyboard = append(keyboard, navRow)
	}

	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙️", "back")))
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func InlineCalendarChanger(year int, month int, readMinutes map[int]int) tgbotapi.InlineKeyboardMarkup {
	daysOfWeek := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var weekRow []tgbotapi.InlineKeyboardButton

	for _, day := range daysOfWeek {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}
	keyboard = append(keyboard, weekRow)

	// Get first weekday and total days in month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	startWeekday := int(firstDay.Weekday()) // Sunday = 0
	if startWeekday == 0 {
		startWeekday = 7 // Adjust for Monday start
	}
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	var row []tgbotapi.InlineKeyboardButton
	// Fill empty slots before the first day
	for i := 1; i < startWeekday; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}

	// Fill in the actual days
	var count int
	for day := 1; day <= daysInMonth; day++ {
		minutes := readMinutes[day]
		if minutes == 0 && day <= time.Now().Day() {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(day)+"❌", fmt.Sprintf("change_%d.%d.%d", day, month, year)))
		} else if minutes != 0 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(day)+"✅", fmt.Sprintf("day_%d.%d.%d", day, month, year)))
			count++
		} else {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(day), fmt.Sprintf("day_%d.%d.%d", day, month, year)))
		}

		// Break at the end of each week
		if len(row) == 7 {
			keyboard = append(keyboard, row)
			row = nil
		}
	}

	// Add remaining row if not complete
	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
		}
		keyboard = append(keyboard, row)
	}

	navRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(count)+"✅", "ignore"),
		tgbotapi.NewInlineKeyboardButtonData("🔄 Update", fmt.Sprintf("back_%d_%d", time.Now().Year(), int(time.Now().Month()))),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(time.Now().Day()-count)+"❌", "ingore"),
	}

	keyboard = append(keyboard, navRow)
	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙️", "back")))

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

//func numberWithEmoji(n int) string {
//	if n < 1 || n > 31 {
//		return "❌ Неверное число"
//	}
//
//	digits := map[rune]string{
//		'0': "0️⃣", '1': "1️⃣", '2': "2️⃣", '3': "3️⃣", '4': "4️⃣",
//		'5': "5️⃣", '6': "6️⃣", '7': "7️⃣", '8': "8️⃣", '9': "9️⃣",
//	}
//
//	var emoji strings.Builder
//	for _, digit := range fmt.Sprintf("%d", n) {
//		emoji.WriteString(digits[digit])
//	}
//	return emoji.String()
//}
