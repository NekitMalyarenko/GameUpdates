package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"strings"
	"os"
)


const (
	BOTTOKEN = "456900455:AAF2uhU9KSd6Gsld4c2M_eZ9b_HDQHggsEI"
)


var bot *tgbotapi.BotAPI


func startBot() {
	var err error

	bot, err = tgbotapi.NewBotAPI(BOTTOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)

		temp := games[PUBG]

		if strings.Contains(update.Message.Text, "/subscribe") {

			if temp.subscribeUser(db, update.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы успешно подписались на обновления по PUBG")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-то пошло не так.Возможно вы уже подписаны на обновление по PUBG?")
				bot.Send(msg)
			}

		} else if strings.Contains(update.Message.Text, "/unsubscribe") {

			if temp.unSubscribeUser(db, update.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы успешно отписались от обновлений по PUBG")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-то пошло не так.Возможно вы не подписаны на обновление по PUBG?")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			bot.Send(msg)
		}
	}

	//devide()
}


func (g *GameData) notifyUsersAboutUpdate(update UpdateData) {
	users := g.getAllUsers(db)

	for _, temp := range users {
		log.Println("notifying user with id(", temp.TelegramId,") about update in", g.GameShortName)
		msg := tgbotapi.NewMessage(int64(temp.TelegramId), update.Url)

		if bot != nil {
			bot.Send(msg)
		} else {
			log.Fatal("!!!!!BOT IS NIL!!!!!")
			os.Exit(228)
		}
	}
}