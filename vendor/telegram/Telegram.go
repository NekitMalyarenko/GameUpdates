package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
	"math/rand"
	"data"
	"db"
	"data/telegram"
	"fmt"
	"data/config"
)


const (
	ITEMS_PER_PAGE = 3
)


var (
	bot *tgbotapi.BotAPI

	/*defaultAnswers = []string{"К сожалению,я не понимаю тебя =(", "А можно по-русски?",
		"Что ты несеш?", "Я тебе не Sire", "Я тебе не Алиса"}*/
	defaultAnswers = []string{"Sorry but i don't understand you =(", "What?",
	"What a hell are you talking about?", "I am not a Siri ", "I am not a Google Assistant"}
)


func StartBot() {
	var err error

	bot, err = tgbotapi.NewBotAPI(configData.GetString(configData.TELEGRAM_BOT_TOKEN))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	telegramData.CreateNextAction()

	updates, err := bot.GetUpdatesChan(u)
	bot.Send(tgbotapi.NewMessage(telegramData.ME, "I'm alive"))

	for update := range updates {
		res := handleMessage(update)

		log.Println(res, len(res))

		if res == nil || len(res) == 0 || res[0] == nil {

			if update.Message != nil {
				_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, defaultAnswers[randInt(0, len(defaultAnswers))]))
			} else {
				_, _ = bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, defaultAnswers[randInt(0, len(defaultAnswers))]))
			}
		}

		for _, val := range res {
			_, err = bot.Send(val)
			if err != nil {
				log.Fatal(res)
				log.Fatal(err)
			}
		}
	}
}


func putOnHold(chatId int64, messageId int){
	holdMessage := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0)})
	bot.Send(holdMessage)
}


func NotifyUsersAboutUpdate(game *data.GameData, message string) {
	users := db.GetDBManager().GetAllUsers(game)

	for _, temp := range users {
		fmt.Println("notifying user with id(", temp.TelegramId,") about update in", game.GameShortName)
		msg := tgbotapi.NewMessage(int64(temp.TelegramId), message)

		if bot != nil {
			bot.Send(msg)
		} else {
			os.Exit(228)
		}
	}
}


func randInt(min int, max int) int {
	return min + rand.Intn(max - min)
}