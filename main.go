package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

const TOKEN = "6892611207:AAHL_s89UFIKt1lTSYnN6OZILcp0zn-D5VU"

var VothmognieVarianti = [3]string{"напиши", "расскажи", "отправь"}
var chatID int64
var result = ""

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

func connectWithTelegram() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		panic("Cannot connect to telegram")
	}
}

func sendMessage(str string) {
	config := tgbotapi.NewMessage(chatID, str)
	bot.Send(config)
}

func isMessageForTelegram(update *tgbotapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	msgInLowerCase := strings.ToLower(update.Message.Text)

	for _, name := range VothmognieVarianti {
		if strings.Contains(msgInLowerCase, name) {
			return true
		}
	}

	return false

}

func getAnswer() string {
	result = ""

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golaing")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		// В телеграм боте айди будут начинаться с 7 так как первые 6 я добавил вручную через phpMyAdmin
		result += fmt.Sprintf("Айди добавленной статьи : %d \nЗaголовок : %s \nАнонс : %s \nТекст : %s \n \n", post.Id, post.Title, post.Anons, post.FullText)
	}

	return result

}

func sendAnswer(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(chatID, getAnswer())
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func main() {
	connectWithTelegram()

	updateConfig := tgbotapi.NewUpdate(0)
	for update := range bot.GetUpdatesChan(updateConfig) {
		if update.Message != nil && update.Message.Text == "/start" {
			chatID = update.Message.Chat.ID
			sendMessage("Этот канал является каналом \"рассылкой\" в котором ты будешь получать информацию с сайта")
		}

		if isMessageForTelegram(&update) {
			sendAnswer(&update)
		}
	}
}
