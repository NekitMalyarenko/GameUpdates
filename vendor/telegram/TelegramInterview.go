package telegram


import (
	"github.com/Syfaro/telegram-bot-api"
	"db"
	"strconv"
	"data/telegram"
	"fmt"
)


const (
	ACTION_START_INTERVIEW         = "start_interview"
	ACTION_CANCEL_INTERVIEW        = "cancel_interview"
	ACTION_NEXT_QUESTION_INTERVIEW = "nextQuestionInter"

	INDEX_QUESTION_NUMBER          = "question_number"

	RATING_QUESTION                = 1
	NO_EXTRA                       = -1
)

var (
	interviewQuestionsData = []string{"Can you rate bot work?", "Which game would you like to add?", "What would you like to add?"}
	interviewExtraData     = []int{RATING_QUESTION, NO_EXTRA, NO_EXTRA}
	star = "\u2b50\ufe0f"
)

/*
	1)Which game would you like to add?
	2)Are you satisfied with bot's work?
	3)What would you like to add?
 */



func StartInterviews() {
	res := make([]MyButtonData, 2)

	callbackData := MyCallbackData{Action : ACTION_START_INTERVIEW}
	res[0] = MyButtonData{Text : "Yes", CallbackData : callbackData.ToJson(), IsNewRow : true}

	callbackData = MyCallbackData{Action : ACTION_CANCEL_INTERVIEW}
	res[1] = MyButtonData{Text : "No", CallbackData : callbackData.ToJson(), IsNewRow : true}

	for _, val := range db.GetDBManager().GetAllUsers(nil) {
		telegramId := val.TelegramId

		msg := tgbotapi.NewMessage(telegramId, "Hello! Could you answer a few questions,for improvement of my work?")
		msg.ReplyMarkup = toInlineKeyboard(res)

		bot.Send(msg)
		fmt.Println("sent interview invitation", telegramId)
	}
}


func getExtra(id int) tgbotapi.InlineKeyboardMarkup {

	switch id {

	case RATING_QUESTION:
		return getRatingBar()

	default:
		return tgbotapi.InlineKeyboardMarkup {
			InlineKeyboard : nil,
		}
	}

}


func getRatingBar() tgbotapi.InlineKeyboardMarkup {
	res := make([]MyButtonData, 5)

	callbackData := MyCallbackData{Action : ACTION_NEXT_QUESTION_INTERVIEW, Temp : "5"}
	res[0] = MyButtonData{Text : star + star + star + star + star, CallbackData : callbackData.ToJson(), IsNewRow : true}

	callbackData = MyCallbackData{Action : ACTION_NEXT_QUESTION_INTERVIEW, Temp : "4"}
	res[1] = MyButtonData{Text : star + star + star + star, CallbackData : callbackData.ToJson(), IsNewRow : true}

	callbackData = MyCallbackData{Action : ACTION_NEXT_QUESTION_INTERVIEW, Temp : "3"}
	res[2] = MyButtonData{Text : star + star + star, CallbackData : callbackData.ToJson(), IsNewRow : true}

	callbackData = MyCallbackData{Action : ACTION_NEXT_QUESTION_INTERVIEW, Temp : "2"}
	res[3] = MyButtonData{Text : star + star, CallbackData : callbackData.ToJson(), IsNewRow : true}

	callbackData = MyCallbackData{Action : ACTION_NEXT_QUESTION_INTERVIEW, Temp : "1"}
	res[4] = MyButtonData{Text : star, CallbackData : callbackData.ToJson(), IsNewRow : true}

	return toInlineKeyboard(res)
}


func getInterviewQuestion(index int) (string, int) {

	if index >= len(interviewQuestionsData) {
		return "", -1
	} else {
		return "(" + strconv.Itoa(index + 1)  + "/" + strconv.Itoa(len(interviewQuestionsData)) + ")\n" + interviewQuestionsData[index], interviewExtraData[index]
	}
}


func sendResultToMe(data *telegramData.InterviewAnswersData) {
	msg := tgbotapi.NewMessage(telegramData.ME, strconv.FormatInt(data.TelegramId, 10) + ":\n" + fmt.Sprintln(data.Answers))
	bot.Send(msg)
	telegramData.RemoveInterviewData(data.TelegramId)
}