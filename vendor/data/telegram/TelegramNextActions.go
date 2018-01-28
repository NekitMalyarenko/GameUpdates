package telegramData


type NextActionData struct {
	Action string
}

var data map[int64]*NextActionData


func Create() {
	data = make(map[int64]*NextActionData)
}


func RegisterData(telegramId int64, action NextActionData) {

	data[telegramId] = &action
}


func UnRegisterData(telegramId int64) {
	data[telegramId] = nil
}


func GetData(telegramId int64) *NextActionData{
	return data[telegramId]
}