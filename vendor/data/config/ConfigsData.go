package configData

import (
	"os"
	"strconv"
	"log"
)

const (
	WEB_PARSER_TIMEOUT      = "GameUpdates_WebParserTimeout"
	CONNECTION_STRING       = "GameUpdates_ConnectionString"
	TELEGRAM_BOT_TOKEN      = "GameUpdates_BotToken"
)


func GetString(varName string) string {
	return os.Getenv(varName)
}


func GetInt(varName string) int {
	temp, err := strconv.Atoi(os.Getenv(varName))
	if err != nil {
		log.Fatal(err)
		return -1
	} else {
		return temp
	}
}


