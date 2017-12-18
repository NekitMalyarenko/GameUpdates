package telegram

import (
	"strings"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"encoding/json"
	"../data"
	"../db"
)


const(
	ACTION_SUBSCRIBE   = "subscribe"
	ACTION_UNSUBSCRIBE = "unsubscribe"
	ACTION_CANCEL      = "cancel"
)


type MyCallbackData struct {
	GameId int	  `json:"game_id"`
	Action string `json:"action"`
}



func handleMessage(update tgbotapi.Update) tgbotapi.Chattable {

	if update.CallbackQuery != nil {
		chatId := update.CallbackQuery.Message.Chat.ID
		messageId := update.CallbackQuery.Message.MessageID
		var callbackData MyCallbackData

		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackData)
		if err != nil {
			log.Fatal(err)
		}

		switch callbackData.Action {

		case ACTION_SUBSCRIBE:
			game := data.GetGame(callbackData.GameId)

			if db.SubscribeUser(game, chatId) {
				return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно подписались на обновления по " + game.GameShortName)
			} else {
				return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы уже подписаны на обновление по " + game.GameShortName + "?")
			}

		case ACTION_UNSUBSCRIBE:
			game := data.GetGame(callbackData.GameId)

			if db.UnSubscribeUser(game, chatId) {
				return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно отписались от обновлений по " + game.GameShortName)
			} else {
				return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы не подписаны на обновление по " + game.GameShortName + "?")
			}


		case ACTION_CANCEL:
			return tgbotapi.DeleteMessageConfig{ChatID : chatId, MessageID : messageId}

		default:
			return nil
		}

	} else {
		chatId := update.Message.Chat.ID
		text := update.Message.Text

		if strings.Contains(text, "/unsubscribe") {
			user, err := db.GetUser(chatId)
			if err != nil {
				log.Fatal(err)
				return tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			}

			keyboard := getUnSubscribeKeyboard(user)
			var msg tgbotapi.MessageConfig

			if keyboard.InlineKeyboard != nil {
				msg = tgbotapi.NewMessage(chatId, "Вы можите отписаться от:")
				msg.ReplyMarkup = keyboard
			} else {
				msg = tgbotapi.NewMessage(chatId, "Вы не подписаны не на одну игру!")
			}

			return msg

		} else if strings.Contains(text, "/subscribe") {
			user, err := db.GetUser(chatId)
			if err != nil {
				log.Fatal(err)
				return tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			}

			keyboard := getSubscribeKeyboard(user)
			var msg tgbotapi.MessageConfig

			if keyboard.InlineKeyboard != nil {
				msg = tgbotapi.NewMessage(chatId, "Вы можите подписаться на:")
				msg.ReplyMarkup = keyboard
			} else {
				msg = tgbotapi.NewMessage(chatId, "Вы уже подписаны на все игры =(")
			}

			return msg

		}/*else if strings.Contains(text, "/cancel") {

		}*/

	}

	return nil
}


func getUnSubscribeKeyboard(u db.User) tgbotapi.InlineKeyboardMarkup {

	res := make([][]tgbotapi.InlineKeyboardButton, 0)
	i := 0

	for _, tempId := range u.Subscribes {
		callbackData, err := json.Marshal(MyCallbackData{Action : ACTION_UNSUBSCRIBE, GameId : tempId})
		if err != nil {
			log.Fatal(err)
		}

		res = append(res, make([]tgbotapi.InlineKeyboardButton, 0))
		res[i] = append(res[i], tgbotapi.NewInlineKeyboardButtonData(data.GetGame(tempId).GameFullName, string(callbackData)))
		i += 1
	}

	res = append(res, getBottomKeyboard())

	if len(res) == 1 {
		return tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: nil,
		}
	} else {
		return tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: res,
		}
	}
}


func getSubscribeKeyboard(u db.User) tgbotapi.InlineKeyboardMarkup {
	res := make([][]tgbotapi.InlineKeyboardButton, 0)
	i := 0

	contains := func(id int) bool{

		for _, value := range u.Subscribes {
			if id == value {
				return true
			}
		}

		return false
	}

	for _, value := range data.GetGames() {

		if !contains(value.GameId) {
			callbackData, err := json.Marshal(MyCallbackData{Action : ACTION_SUBSCRIBE, GameId : value.GameId})
			if err != nil {
				log.Fatal(err)
			}

			res = append(res, make([]tgbotapi.InlineKeyboardButton, 0))
			res[i] = append(res[i], tgbotapi.NewInlineKeyboardButtonData(value.GameFullName, string(callbackData)))
			i += 1
		}
	}

	res = append(res, getBottomKeyboard())

	if len(res) == 1 {
		return tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: nil,
		}
	} else {
		return tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: res,
		}
	}
}


func getBottomKeyboard() []tgbotapi.InlineKeyboardButton{
	result := make([]tgbotapi.InlineKeyboardButton, 0)
	callbackData := MyCallbackData{Action : ACTION_CANCEL, GameId : -1}

	result = append(result,  tgbotapi.NewInlineKeyboardButtonData("Отмена", callbackData.toJson()))

	return result
}


func (callbackData *MyCallbackData) toJson() string {
	res, err := json.Marshal(callbackData)
	if err != nil {
		log.Fatal(err)
	}

	return string(res)
}


func (callbackData *MyCallbackData) fromJson(input []byte) {
	err := json.Unmarshal(input, &callbackData)
	if err != nil {
		log.Fatal(err)
	}
}


/*func getGameId(inputString string) int {
	startIndex := strings.LastIndex(inputString, ":") + 1
	runes := []byte(inputString)

	gameId, err := strconv.Atoi(string(runes[startIndex:]))
	if err != nil {
		log.Fatal(err)
	}

	return gameId
}


func autoComplete(inputString string, update tgbotapi.Update) string {

	res := inputString

	res = strings.Replace(inputString, "@bot_name", bot.Self.FirstName + " " + bot.Self.LastName, -1)
	res = strings.Replace(inputString, "@user_name", update.Message.From.FirstName + " " + update.Message.From.LastName, -1)

	return res
}*/