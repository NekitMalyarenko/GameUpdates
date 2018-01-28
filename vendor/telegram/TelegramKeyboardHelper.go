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
	/*var firstRow = []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Все игры"), tgbotapi.NewKeyboardButton("Мои подписки")}
	var secondRow = []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Поиск")}*/
	var firstRow = []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("All Games"), tgbotapi.NewKeyboardButton("My Subscribes")}
	var secondRow = []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("Search")}

	return tgbotapi.NewReplyKeyboard(firstRow, secondRow)
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
		res = append(res, MyButtonData{Text : "\u2b05\ufe0f", CallbackData : callbackData.ToJson(), IsNewRow : true})
	}

	if hasNext {
		callbackData = MyCallbackData{Action : ACTION_CHANGE_PAGE, Page : page + 1, Temp : pageAction}

		if page == 0 {
			res = append(res, MyButtonData{Text : "\u27a1\ufe0f", CallbackData : callbackData.ToJson(), IsNewRow : true})
		} else {
			res = append(res, MyButtonData{Text : "\u27a1\ufe0f", CallbackData : callbackData.ToJson(), IsNewRow : false})
		}
	}

	return res
}


func getCancelButton() MyButtonData {
	callbackData := MyCallbackData{Action : ACTION_CANCEL}

	return MyButtonData{
		Text : "\u274c",
		CallbackData : callbackData.ToJson(),
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


func (callbackData *MyCallbackData) ToJson() string {
	res, err := json.Marshal(callbackData)
	if err != nil {
		log.Fatal(err)
	}

	return string(res)
}