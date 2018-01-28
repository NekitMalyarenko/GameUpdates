package telegramData


type NextActionData struct {
	Action   string
	TempData map[string]interface{}
}

var data map[int64]*NextActionData


func CreateNextAction() {
	data = make(map[int64]*NextActionData)
}

func RegisterNextActionData(telegramId int64, action NextActionData) {
	data[telegramId] = &action
}

func UnRegisterNextActionData(telegramId int64) {
	data[telegramId] = nil
}

func GetNextActionData(telegramId int64) *NextActionData{
	return data[telegramId]
}