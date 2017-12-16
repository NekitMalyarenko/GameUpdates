package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
	"math/rand"
	"GameUpdates/db"
	"GameUpdates/data"
)


const (
	BOTTOKEN = "456900455:AAF2uhU9KSd6Gsld4c2M_eZ9b_HDQHggsEI"
)


var bot *tgbotapi.BotAPI

var defaultAnswers = []string{"К сожалению,я не понимаю тебя =(", "А можно по-русски?",
								"Что ты несеш?", "Я тебе не Sire", "Я тебе не Алиса"}


func StartBot() {
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
		//log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)

		//temp := games[PUBG]

		res := handleMessage(update)

		if res == nil {

			if update.Message != nil {
				res = tgbotapi.NewMessage(update.Message.Chat.ID, defaultAnswers[randInt(0, len(defaultAnswers))])
			} else {
				res = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, defaultAnswers[randInt(0, len(defaultAnswers))])
			}
		}

		_, err = bot.Send(res)
		if err != nil {
				log.Fatal(res)
				log.Fatal(err)
			}
	}

	//devide()
}


func NotifyUsersAboutUpdate(game *data.GameData, message string) {
	users := db.GetAllUsers(game)

	for _, temp := range users {
		log.Println("notifying user with id(", temp.TelegramId,") about update in", game.GameShortName)
		msg := tgbotapi.NewMessage(int64(temp.TelegramId), message)

		if bot != nil {
			bot.Send(msg)
		} else {
			log.Fatal("!!!!!BOT IS NIL!!!!!")
			os.Exit(228)
		}
	}
}


func randInt(min int, max int) int {
	return min + rand.Intn(max - min)
}