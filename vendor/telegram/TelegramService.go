package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"data"
	"db"
	"data/telegram"
)

const(
	ACTION_UNSUBSCRIBE = "action_1"
	ACTION_SUBSCRIBE   = "subscribe"
	ACTION_CLICK       = "click"
	ACTION_SEARCH      = "search"

	ACTION_CANCEL      = "cancel"

	ACTION_CHANGE_PAGE = "page_change"
	ACTION_FIRST_PAGE  = "page_first"
	ACTION_LAST_PAGE   = "page_last"

	GREETING_MESSAGE   = "Привет,для того чтобы оповещать тебя об обновлениях в играх,я должен знать что тебе интересно." +
		"Чтобы подписаться на обновление по игре:Нажми кнопку 'Search' и ввиди название игры,после чего выбери ее и нажми 'Подписаться'."
)


func handleMessage(update tgbotapi.Update) tgbotapi.Chattable {

	if update.CallbackQuery != nil {
		return handleCallbackQuery(update)
	} else {
		return handleText(update)
	}
}


func handleCallbackQuery(update tgbotapi.Update) tgbotapi.Chattable{
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	go putOnHold(chatId, messageId)
	callbackData := fromJson([]byte(update.CallbackQuery.Data))

	switch callbackData.Action {

	case ACTION_SUBSCRIBE:
		game := data.GetGame(callbackData.GameId)

		if db.GetDBManager().SubscribeUser(game, chatId) {
			return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно подписались на обновления по " + game.GameShortName)
		} else {
			return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы уже подписаны на обновление по " + game.GameShortName + "?")
		}
		break

	case ACTION_UNSUBSCRIBE:
		game := data.GetGame(callbackData.GameId)

		if db.GetDBManager().UnSubscribeUser(game, chatId) {
			return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно отписались от обновлений по " + game.GameShortName)
		} else {
			return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы не подписаны на обновление по " + game.GameShortName + "?")
		}
		break

	case ACTION_CHANGE_PAGE, ACTION_LAST_PAGE, ACTION_FIRST_PAGE:
		page := callbackData.Page
		pageAction := callbackData.Temp
		user, err := db.GetDBManager().GetUser(chatId)
		if err != nil {
			log.Println(err)
			return tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		}

		switch pageAction {

		case ACTION_SUBSCRIBE:
			return tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, getSubscribeKeyboard(user, page))

		case ACTION_UNSUBSCRIBE:
			return tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, getUnSubscribeKeyboard(user, page))
		}

		break

	case ACTION_CANCEL:
		return tgbotapi.DeleteMessageConfig{ChatID : chatId, MessageID : messageId}

	case ACTION_CLICK:
		isSubscribed := false
		user, err := db.GetDBManager().GetUser(chatId)
		if err != nil {
			log.Fatal(err)
			return tgbotapi.NewMessage(chatId, err.Error())
		}

		for _, gameId := range user.Subscribes {

			if gameId == callbackData.GameId {
				isSubscribed = true
				break
			}
		}

		res := make([]MyButtonData, 0)

		if isSubscribed {
			callbackData := MyCallbackData{Action : ACTION_UNSUBSCRIBE, GameId : callbackData.GameId}
			res = append(res, MyButtonData{Text : "Отписаться от " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
		} else {
			callbackData := MyCallbackData{Action : ACTION_SUBSCRIBE, GameId : callbackData.GameId}
			res = append(res, MyButtonData{Text : "Подписаться на " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
		}

		msg := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, toInlineKeyboard(res))

		return msg

	}

	return nil
}


func handleText(update tgbotapi.Update) tgbotapi.Chattable {
	chatId := update.Message.Chat.ID
	text := update.Message.Text
	user, err := db.GetDBManager().GetUser(chatId)
	if err != nil {
		log.Fatal(err)
	}

	var msg tgbotapi.Chattable = nil

	switch text {

	case "Search":
		telegramData.RegisterData(chatId, telegramData.NextActionData{Action : ACTION_SEARCH})
		msg = tgbotapi.NewMessage(chatId, "Напишите примерно как называется игра:")
		break


	case "My Subscribes":
		keyboard := getUnSubscribeKeyboard(user, 0)

		if keyboard.InlineKeyboard != nil {
			temp := tgbotapi.NewMessage(chatId, "Нажмите на игру чтобы отписаться от нее:")
			temp.ReplyMarkup = keyboard
			msg = temp
		} else {
			msg = tgbotapi.NewMessage(chatId, "Вы не подписаны не на одну игру,нажмите на Search и найдите интересующую вас игру.")
		}

		break

	case "/start":
		telegramData.UnRegisterData(chatId)
		temp := tgbotapi.NewMessage(chatId, GREETING_MESSAGE)
		temp.ReplyMarkup = getKeyboard()
		msg = temp
		break
	}


	if msg == nil && telegramData.GetData(chatId) != nil {
		msg = handleNextAction(update, chatId, text)
	}

	return msg
}


func handleNextAction(update tgbotapi.Update, chatId int64, text string) tgbotapi.Chattable {
	var msg tgbotapi.Chattable = nil
	nextAction := telegramData.GetData(chatId)
	log.Println("nextAction:", nextAction)

	switch nextAction.Action {

	case ACTION_SEARCH:
		keyboard := getSearchResultsKeyboard(text)

		if keyboard.InlineKeyboard == nil {
			msg = tgbotapi.NewMessage(chatId, "К сожелению,я ничего не нашел.")
		} else {
			temp := tgbotapi.NewMessage(chatId, "Вот что я нашел:")
			temp.ReplyMarkup = keyboard
			msg = temp
		}

		telegramData.UnRegisterData(chatId)
		break
	}


	return msg
}


func getUnSubscribeKeyboard(u db.User, page int) tgbotapi.InlineKeyboardMarkup {

	if len(u.Subscribes) != 0 {

		res := make([]MyButtonData, 0)

		startIndex := page * ITEMS_PER_PAGE
		endIndex := startIndex + ITEMS_PER_PAGE
		temp := 0

		for i := 0; i < len(u.Subscribes); i++ {

			if startIndex <= i && endIndex > i {
				callbackData := MyCallbackData{Action: ACTION_UNSUBSCRIBE, GameId: i}
				res = append(res, MyButtonData{Text: data.GetGame(u.Subscribes[temp]).GameFullName, CallbackData: callbackData.toJson(), IsNewRow: true})
				temp++
			}
		}

		if len(u.Subscribes) <= endIndex {
			res = getBottom(res, false, page)
		} else {
			res = getBottom(res, true, page)
		}

		return toInlineKeyboard(res)

	} else {
		return tgbotapi.InlineKeyboardMarkup {
			InlineKeyboard : nil,
		}
	}
}


func getSearchResultsKeyboard(request string) tgbotapi.InlineKeyboardMarkup {
	games := db.GetDBManager().SearchGame(request)

	if len(games) != 0 {
		res := make([]MyButtonData, 0)

		for _, gameId := range games {
			callbackData := MyCallbackData{Action : ACTION_CLICK, GameId : gameId}
			res = append(res, MyButtonData{Text : data.GetGame(gameId).GameFullName, CallbackData : callbackData.toJson(), IsNewRow : true})
		}

		res = append(res, getCancelButton())
		return toInlineKeyboard(res)

	} else {
		return tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard : nil,
		}
	}
}


func getSubscribeKeyboard(u db.User, page int) tgbotapi.InlineKeyboardMarkup {
	res := make([]MyButtonData, 0)

	games := getUniqueGamesForUser(u)
	startIndex := page * ITEMS_PER_PAGE
	endIndex := startIndex + ITEMS_PER_PAGE

	if len(games) != 0 {

		var callbackData MyCallbackData

		for i := 0; i < len(games); i++ {

			log.Println(startIndex, "<=", i, "<", endIndex, startIndex <= i && endIndex >= i)

			if startIndex <= i && endIndex > i {
				callbackData = MyCallbackData{Action: ACTION_SUBSCRIBE, GameId: i}
				res = append(res, MyButtonData{Text: data.GetGame(games[i]).GameFullName, CallbackData: callbackData.toJson(), IsNewRow: true})
			}

		}

		if len(games) <= endIndex {
			getBottom(res, false, page)
		} else {
			getBottom(res, true, page)
		}

		return toInlineKeyboard(res)

	} else {
		return tgbotapi.InlineKeyboardMarkup {
			InlineKeyboard : nil,
		}
	}
}


func getBottom(data []MyButtonData, hasNext bool, page int) []MyButtonData{
	nav := getNavigationButtons(page, hasNext, ACTION_SUBSCRIBE)

	for _, value := range nav {
		data = append(data, value)
	}

	data = append(data, getCancelButton())

	return data
}


func getUniqueGamesForUser(u db.User) []int {
	res := make([]int, 0)
	contains := func(id int) bool{
		for _, value := range u.Subscribes {
			if id == value {
				return true
			}
		}

		return false
	}

	for i :=1; i < len(data.GetGames()) + 1; i++ {
		log.Println("gameId:", i, "contains:", contains(i))

		if !contains(i) {
			res = append(res, i)
		}
	}

	return res
}


func join(to [][]tgbotapi.InlineKeyboardButton, from [][]tgbotapi.InlineKeyboardButton) [][]tgbotapi.InlineKeyboardButton {

	for index := range from {
		to = append(to, make([]tgbotapi.InlineKeyboardButton, 0))

		for _, element := range from[index] {
			to[len(to) - 1] = append(to[len(to) - 1], element)
		}

	}

	return to
}