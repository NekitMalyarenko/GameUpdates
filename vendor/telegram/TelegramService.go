package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"data"
	"db"
	"data/telegram"
	"strconv"
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


func handleMessage(update tgbotapi.Update) []tgbotapi.Chattable {

	if update.CallbackQuery != nil {
		return handleCallbackQuery(update)
	} else {
		return handleText(update)
	}
}


func handleCallbackQuery(update tgbotapi.Update) []tgbotapi.Chattable {
	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	go putOnHold(chatId, messageId)
	callbackData := fromJson([]byte(update.CallbackQuery.Data))

	res := make([]tgbotapi.Chattable, 0)

	switch callbackData.Action {

	case ACTION_SUBSCRIBE:
		game := data.GetGame(callbackData.GameId)

		if db.GetDBManager().SubscribeUser(game, chatId) {
			//return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно подписались на обновления по " + game.GameShortName)
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, "You successfully subscribed on " + game.GameShortName))
		} else {
			//return tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы уже подписаны на обновление по " + game.GameShortName + "?")
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId,"Smth went wrong,sorry =("))
		}
	break

	case ACTION_UNSUBSCRIBE:
		log.Println("callbackData:", callbackData)
		game := data.GetGame(callbackData.GameId)

		if db.GetDBManager().UnSubscribeUser(game, chatId) {
			//return tgbotapi.NewEditMessageText(chatId, messageId, "Вы успешно отписались от обновлений по " + game.GameShortName)
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, "You successfully unsubscribed from " + game.GameShortName))
		} else {
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId,"Что-то пошло не так,возможно вы не подписаны на обновление по " + game.GameShortName + "?"))
		}
	break

	case ACTION_CHANGE_PAGE, ACTION_LAST_PAGE, ACTION_FIRST_PAGE:
		page := callbackData.Page
		pageAction := callbackData.Temp
		user, err := db.GetDBManager().GetUser(chatId)
		if err != nil {
			log.Println(err)
			res = append(res, tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		}

		switch pageAction {

		case ACTION_SUBSCRIBE:
			res = append(res, tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, getSubscribeKeyboard(user, page)))

		case ACTION_UNSUBSCRIBE:
			res = append(res, tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, getUnSubscribeKeyboard(user, page)))
		}

	break

	case ACTION_CANCEL:
		res = append(res, tgbotapi.DeleteMessageConfig{ChatID : chatId, MessageID : messageId})
	break

	case ACTION_CLICK:
		isSubscribed := false
		user, err := db.GetDBManager().GetUser(chatId)
		if err != nil {
			log.Fatal(err)
			res = append(res, tgbotapi.NewMessage(chatId, err.Error()))
		}

		for _, gameId := range user.Subscribes {

			if gameId == callbackData.GameId {
				isSubscribed = true
				break
			}
		}

		result := make([]MyButtonData, 0)

		if isSubscribed {
			callbackData := MyCallbackData{Action : ACTION_UNSUBSCRIBE, GameId : callbackData.GameId}
			//res = append(res, MyButtonData{Text : "ОТПИСАТЬСЯ от " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.toJson()})
			result = append(result, MyButtonData{Text : "UNSUBSCRIBE from " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.ToJson()})
		} else {
			callbackData := MyCallbackData{Action : ACTION_SUBSCRIBE, GameId : callbackData.GameId}
			result = append(result, MyButtonData{Text : "SUBSCRIBE on " + data.GetGame(callbackData.GameId).GameShortName, CallbackData : callbackData.ToJson()})
		}

		result = append(result, getCancelButton())
		msg := tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, toInlineKeyboard(result))

		res = append(res, msg)
		break

	case ACTION_START_INTERVIEW:
		telegramData.UnRegisterNextActionData(chatId)

		tempData := make(map[string]interface{}, 0)
		tempData[INDEX_QUESTION_NUMBER] = 0

		nextAction := telegramData.NextActionData{
			Action : ACTION_NEXT_QUESTION_INTERVIEW,
			TempData : tempData,
		}

		telegramData.RegisterNextActionData(chatId, nextAction)
		telegramData.CreateAnswerData(chatId)

		question, extra := getInterviewQuestion(0)

		if extra == NO_EXTRA {
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, question))
		} else {
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, question))
			res = append(res, tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, getExtra(extra)))
		}

		break

	case ACTION_NEXT_QUESTION_INTERVIEW:
		nextAction := telegramData.GetNextActionData(chatId)
		id := (nextAction.TempData[INDEX_QUESTION_NUMBER]).(int) + 1

		question, extra := getInterviewQuestion(id)
		telegramData.AddAnswer(chatId, callbackData.Temp)

		if question != "" {
			nextAction.TempData[INDEX_QUESTION_NUMBER] = id

			msg := tgbotapi.NewMessage(chatId, question)

			if extra != NO_EXTRA {
				msg.ReplyMarkup = getExtra(extra)
			}

			res = append(res, msg)
			res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, "You rated " + callbackData.Temp))
		} else {
			telegramData.UnRegisterNextActionData(chatId)
			res = append(res, tgbotapi.NewMessage(chatId, "Thank you for your answers =)"))
			sendResultToMe(telegramData.GetInterviewData(chatId))
		}
		break

	case ACTION_CANCEL_INTERVIEW:
		res = append(res, tgbotapi.NewEditMessageText(chatId, messageId, "=("))
		res = append(res, tgbotapi.NewMessage(360952996, strconv.FormatInt(chatId, 10) + " canceled interview"))
	break
	}

	return res
}


func handleText(update tgbotapi.Update) []tgbotapi.Chattable {
	chatId := update.Message.Chat.ID
	text := update.Message.Text
	user, err := db.GetDBManager().GetUser(chatId)
	if err != nil {
		log.Fatal(err)
	}
	response := make([]tgbotapi.Chattable, 0)

	switch text {

	case "Search":
		telegramData.RegisterNextActionData(chatId, telegramData.NextActionData{Action : ACTION_SEARCH})
		//msg = tgbotapi.NewMessage(chatId, "Напишите примерно как называется игра:")
		response = append(response, tgbotapi.NewMessage(chatId, "Type name of the game:"))
		break

	case "My Subscribes":
		keyboard := getUnSubscribeKeyboard(user, 0)

		if keyboard.InlineKeyboard != nil {
			//temp := tgbotapi.NewMessage(chatId, "Нажмите на игру чтобы ОТПИСАТЬСЯ от нее:")
			msg := tgbotapi.NewMessage(chatId, "Click on the game to UNSUBSCRIBE from it:")
			msg.ReplyMarkup = keyboard
			response = append(response, msg)
		} else {
			//msg = tgbotapi.NewMessage(chatId, "Вы не подписаны не на одну игру,нажмите на Search и найдите интересующую вас игру.")
			msg := tgbotapi.NewMessage(chatId, "You aren't subscribed on the any game.")
			response = append(response, msg)
		}

		break

	case "All Games":
		keyboard := getSubscribeKeyboard(user, 0)

		if keyboard.InlineKeyboard != nil {
			//temp := tgbotapi.NewMessage(chatId, "Нажмите на игру чтобы ПОДПИСАТЬСЯ на нее:")
			msg := tgbotapi.NewMessage(chatId, "Click on the game to SUBSCRIBE on it:")
			msg.ReplyMarkup = keyboard
			response = append(response, msg)
		}else {
			//msg = tgbotapi.NewMessage(chatId, "Вы уже подписаны на все игры.")
			msg := tgbotapi.NewMessage(chatId, "You are already subscribed on all games.")
			response = append(response, msg)
		}

		break

	case "/start":
		telegramData.UnRegisterNextActionData(chatId)
		msg := tgbotapi.NewMessage(chatId, GREETING_MESSAGE)
		msg.ReplyMarkup = getKeyboard()
		response = append(response, msg)
		break


	/*case "/test":
		res := make([]MyButtonData, 2)

		callbackData := MyCallbackData{Action : ACTION_START_INTERVIEW}
		res[0] = MyButtonData{Text : "Yes", CallbackData : callbackData.ToJson(), IsNewRow : true}

		callbackData = MyCallbackData{Action : ACTION_CANCEL_INTERVIEW}
		res[1] = MyButtonData{Text : "No", CallbackData : callbackData.ToJson(), IsNewRow : false}

		msg := tgbotapi.NewMessage(chatId, "Hello! Could you answer a few questions,for improvement of my work?")
		msg.ReplyMarkup = toInlineKeyboard(res)
		response = append(response, msg)
		break*/

	case "/startInterview":
		if chatId == telegramData.ME {
			go StartInterviews()
			response = append(response, tgbotapi.NewMessage(telegramData.ME, "OK"))
		}
		break
	}


	if len(response) == 0 && telegramData.GetNextActionData(chatId) != nil {
		msg := handleNextAction(update, chatId, text)

		if msg == nil {
			return []tgbotapi.Chattable{}
		} else {
			return []tgbotapi.Chattable{msg}
		}

	}


	return response
}


func handleNextAction(update tgbotapi.Update, chatId int64, text string) tgbotapi.Chattable {
	var msg tgbotapi.Chattable = nil
	nextAction := telegramData.GetNextActionData(chatId)
	log.Println("nextAction:", nextAction)

	switch nextAction.Action {

	case ACTION_SEARCH:
		keyboard := getSearchResultsKeyboard(text)

		if keyboard.InlineKeyboard == nil {
			//msg = tgbotapi.NewMessage(chatId, "К сожелению,я ничего не нашел.")
			msg = tgbotapi.NewMessage(chatId, "Sorry,but i haven't found nothing.")
		} else {
			//temp := tgbotapi.NewMessage(chatId, "Вот что я нашел:")
			temp := tgbotapi.NewMessage(chatId, "That's what i have found:")
			temp.ReplyMarkup = keyboard
			msg = temp
		}

		telegramData.UnRegisterNextActionData(chatId)
		break

	case ACTION_NEXT_QUESTION_INTERVIEW:
		id := (nextAction.TempData[INDEX_QUESTION_NUMBER]).(int) + 1

		question, extra := getInterviewQuestion(id)
		telegramData.AddAnswer(chatId, text)

		if question != "" {
			nextAction.TempData[INDEX_QUESTION_NUMBER] = id
			msg := tgbotapi.NewMessage(chatId, question)

			if extra != NO_EXTRA {
				msg.ReplyMarkup = getExtra(extra)
			}

			return msg
		} else {
			telegramData.UnRegisterNextActionData(chatId)
			sendResultToMe(telegramData.GetInterviewData(chatId))
			return tgbotapi.NewMessage(chatId, "Thank you for your answers =)")
		}
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
				res = append(res, MyButtonData{Text: data.GetGame(u.Subscribes[i]).GameFullName, CallbackData: callbackData.ToJson(), IsNewRow: true})
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
			res = append(res, MyButtonData{Text : data.GetGame(gameId).GameFullName, CallbackData : callbackData.ToJson(), IsNewRow : true})
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
				res = append(res, MyButtonData{Text: data.GetGame(games[i]).GameFullName, CallbackData: callbackData.ToJson(), IsNewRow: true})
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