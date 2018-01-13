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

	/*GREETING_MESSAGE   = "Привет,для того чтобы оповещать тебя об обновлениях в играх,я должен знать что тебе интересно." +
		"Чтобы подписаться на обновление по игре:\n1)Нажми кнопку 'Поиск' и ввиди название игры,после чего выбери ее и нажми 'Подписаться'." +
		"\n2)Нажми 'Все игры' и можешь посмотреть все игры которые я поддерживаю."*/
	GREETING_MESSAGE   = "Hello,I can't alert you about updates,if i don't know what are you interested in." +
	"To subscribe on alerts:\n1)Click on 'Search' and type name of game,after that choose one and click 'Subscribe'." +
	"\n2)Click 'All games' and that's all game which i support.\nIf u have any questions/suggestions write me http://www.t.me/Zikim"
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
			//return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно подписались на обновления по " + game.GameShortName)
			return tgbotapi.NewEditMessageText(chatId, messageId, "You successfully subscribed on " + game.GameShortName)
		} else {
			//return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы уже подписаны на обновление по " + game.GameShortName + "?")
			return tgbotapi.NewEditMessageText(chatId, messageId,"Smth went wrong,sorry =(")
		}
		break

	case ACTION_UNSUBSCRIBE:
		log.Println("callbackData:", callbackData)
		game := data.GetGame(callbackData.GameId)

		if db.GetDBManager().UnSubscribeUser(game, chatId) {
			//return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно отписались от обновлений по " + game.GameShortName)
			return tgbotapi.NewEditMessageText(chatId, messageId, "You successfully unsubscribed from " + game.GameShortName)
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
			//res = append(res, MyButtonData{Text : "ОТПИСАТЬСЯ от " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
			res = append(res, MyButtonData{Text : "UNSUBSCRIBE from " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
		} else {
			callbackData := MyCallbackData{Action : ACTION_SUBSCRIBE, GameId : callbackData.GameId}
			res = append(res, MyButtonData{Text : "SUBSCRIBE on " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
		}

		res = append(res, getCancelButton())
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
		//msg = tgbotapi.NewMessage(chatId, "Напишите примерно как называется игра:")
		msg = tgbotapi.NewMessage(chatId, "Type name of the game:")
		break


	case "My Subscribes":
		keyboard := getUnSubscribeKeyboard(user, 0)

		if keyboard.InlineKeyboard != nil {
			//temp := tgbotapi.NewMessage(chatId, "Нажмите на игру чтобы ОТПИСАТЬСЯ от нее:")
			temp := tgbotapi.NewMessage(chatId, "Click on the game to UNSUBSCRIBE from it:")
			temp.ReplyMarkup = keyboard
			msg = temp
		} else {
			//msg = tgbotapi.NewMessage(chatId, "Вы не подписаны не на одну игру,нажмите на Search и найдите интересующую вас игру.")
			msg = tgbotapi.NewMessage(chatId, "You aren't subscribed on the any game.")
		}

		break

	case "All Games":
		keyboard := getSubscribeKeyboard(user, 0)

		if keyboard.InlineKeyboard != nil {
			//temp := tgbotapi.NewMessage(chatId, "Нажмите на игру чтобы ПОДПИСАТЬСЯ на нее:")
			temp := tgbotapi.NewMessage(chatId, "Click on the game to SUBSCRIBE on it:")
			temp.ReplyMarkup = keyboard
			msg = temp
		}else {
			//msg = tgbotapi.NewMessage(chatId, "Вы уже подписаны на все игры.")
			msg = tgbotapi.NewMessage(chatId, "You are already subscribed on all games.")
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
			//msg = tgbotapi.NewMessage(chatId, "К сожелению,я ничего не нашел.")
			msg = tgbotapi.NewMessage(chatId, "Sorry,but i haven found nothing.")
		} else {
			//temp := tgbotapi.NewMessage(chatId, "Вот что я нашел:")
			temp := tgbotapi.NewMessage(chatId, "That's what i found:")
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
		log.Println("startIndex:", startIndex, "endIndex:", endIndex)

		for i := 0; i < len(u.Subscribes); i++ {

			log.Println( startIndex, "<", i, "<", endIndex)

			if startIndex <= i && endIndex > i {
				callbackData := MyCallbackData{Action: ACTION_CLICK, GameId: u.Subscribes[i]}
				res = append(res, MyButtonData{Text: data.GetGame(u.Subscribes[i]).GameFullName, CallbackData: callbackData.toJson(), IsNewRow: true})
			}
		}

		if len(u.Subscribes) <= endIndex {
			res = getBottom(res, false, page, ACTION_UNSUBSCRIBE)
		} else {
			res = getBottom(res, true, page, ACTION_UNSUBSCRIBE)
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

	log.Println("games:", games)

	if len(games) != 0 {

		var callbackData MyCallbackData

		for i := 0; i < len(games); i++ {

			log.Println(startIndex, "<=", i, "<", endIndex, startIndex <= i && endIndex >= i)

			if startIndex <= i && endIndex > i {
				callbackData = MyCallbackData{Action: ACTION_SUBSCRIBE, GameId: games[i]}
				res = append(res, MyButtonData{Text: data.GetGame(games[i]).GameFullName, CallbackData: callbackData.toJson(), IsNewRow: true})
			}

		}

		if len(games) <= endIndex {
			res = getBottom(res, false, page, ACTION_SUBSCRIBE)
		} else {
			res = getBottom(res, true, page, ACTION_SUBSCRIBE)
		}

		return toInlineKeyboard(res)

	} else {
		return tgbotapi.InlineKeyboardMarkup {
			InlineKeyboard : nil,
		}
	}
}


func getBottom(data []MyButtonData, hasNext bool, page int, action string) []MyButtonData{
	nav := getNavigationButtons(page, hasNext, action)

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