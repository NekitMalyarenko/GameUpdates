package telegram

import (
	"github.com/Syfaro/telegram-bot-api"
	"encoding/json"
	"log"
)

type MyButtonData struct {
	Text         string
	CallbackData string
	IsNewRow       bool
}

type MyCallbackData struct {
	GameId     int	  `json:"game_id"`
	Action     string `json:"action"`
	Page       int    `json:"page"`
	Temp       string `json:"temp"`
}


func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	res := make([]tgbotapi.KeyboardButton, 0)

	res = append(res, tgbotapi.NewKeyboardButton("Search"))
	res = append(res, tgbotapi.NewKeyboardButton("My Subscribes"))

	return tgbotapi.NewReplyKeyboard(res)
}


func toInlineKeyboard(data []MyButtonData) tgbotapi.InlineKeyboardMarkup{
	res := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i := 0; i < len(data); i++ {
		buttonData := data[i]

		if buttonData.IsNewRow {
			res = append(res, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonData.Text, buttonData.CallbackData)))
		} else {

			if len(res) != 0 {
				lastRowId := len(res) - 1
				res[lastRowId] = append(res[lastRowId], tgbotapi.NewInlineKeyboardButtonData(buttonData.Text, buttonData.CallbackData))
			} else {
				res = append(res, make([]tgbotapi.InlineKeyboardButton, 0))
				res[0] = append(res[0], tgbotapi.NewInlineKeyboardButtonData(buttonData.Text, buttonData.CallbackData))
			}
		}

	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard : res,
	}
}


func getNavigationButtons(page int, hasNext bool, pageAction string) []MyButtonData {
	res := make([]MyButtonData, 0)

	var callbackData MyCallbackData

	if page != 0 {
		callbackData = MyCallbackData{Action : ACTION_CHANGE_PAGE, Page : page - 1, Temp : pageAction}
		res = append(res, MyButtonData{Text : "<", CallbackData : callbackData.toJson(), IsNewRow : true})
	}

	if hasNext {
		callbackData = MyCallbackData{Action : ACTION_CHANGE_PAGE, Page : page + 1, Temp : pageAction}

		if page == 0 {
			res = append(res, MyButtonData{Text : ">", CallbackData : callbackData.toJson(), IsNewRow : true})
		} else {
			res = append(res, MyButtonData{Text : ">", CallbackData : callbackData.toJson(), IsNewRow : false})
		}
	}

	return res
}


func getCancelButton() MyButtonData {
	callbackData := MyCallbackData{Action : ACTION_CANCEL}

	return MyButtonData{
		Text : "Закрыть",
		CallbackData : callbackData.toJson(),
		IsNewRow : true,
	}
}


func fromJson(input []byte) MyCallbackData {
	var callbackData MyCallbackData
	err := json.Unmarshal(input, &callbackData)
	if err != nil {
		log.Fatal(err)
	}

	return callbackData
}


func (callbackData *MyCallbackData) toJson() string {
	res, err := json.Marshal(callbackData)
	if err != nil {
		log.Fatal(err)
	}

	return string(res)
}