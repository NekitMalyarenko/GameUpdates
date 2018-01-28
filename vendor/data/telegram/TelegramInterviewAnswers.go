package telegramData

import (
	"sync"
	"log"
)


type InterviewAnswersData struct {
	TelegramId int64
	Answers    []string
}


var (
	once sync.Once
	answers map[int64]*InterviewAnswersData
)


func init() {
	once.Do(func() {
		answers = make(map[int64]*InterviewAnswersData, 0)
		log.Println("initing answers")
	})
}


func CreateAnswerData(telegramId int64) {
	data := InterviewAnswersData{
		TelegramId : telegramId,
		Answers : make([]string, 0),
	}

	answers[telegramId] = &data
}


func AddAnswer(telegramId int64, answer string) bool {
	interviewData, ok := answers[telegramId]

	if ok {
		interviewData.Answers = append(interviewData.Answers, answer)
		return true
	} else {
		return false
	}
}


func GetInterviewData(telegramId int64) *InterviewAnswersData {
	return answers[telegramId]
}


func RemoveInterviewData(telegramId int64) {
	delete(answers, telegramId)
}