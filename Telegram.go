package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"strconv"
	"Syfaro/telegram-bot-api"
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

		/*if strings.Contains(update.Message.Text, "/subscribe") {

			addUser(PUBG, update.Message.From.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы успешно подписались на обновления по PUBG")
			bot.Send(msg)
		}*/

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		bot.Send(msg)
	}

	devide()
}


func (g *GameData) notifyUsersAboutUpdate(update UpdateData) {
	users := getUsers(g.GameId)

	for _, temp := range users {
		parsed, _ := strconv.ParseInt(temp, 10 ,0)
		log.Println("notifying user with id:", parsed,"about update in", g.GameName)
		msg := tgbotapi.NewMessage(parsed, update.Url)
		bot.Send(msg)
	}
}