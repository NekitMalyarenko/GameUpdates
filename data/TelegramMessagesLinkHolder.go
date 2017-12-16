package data

import "time"

type TelegramMessageLink struct{
	TelegramId    int64
	LastMessageId int
	Time          time.Time
}

var links map[int64]*TelegramMessageLink = nil


func AddMessageLink(link TelegramMessageLink){
	links[link.TelegramId] = &link
}


func GetMessagesLinks() map[int64]*TelegramMessageLink {
	return links
}


func RemoveMessage(telegramId int64) {
	delete(links, telegramId)
}