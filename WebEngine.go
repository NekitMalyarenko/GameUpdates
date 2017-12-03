package main

import (
	"github.com/opesun/goquery"
	"strings"
	"time"
	"log"
)

const(
	SLEEPING = 60
)


type UpdateData struct {
	Id  string
	Url string
}

type GameData struct {
	GameId       int    `json:"game_id"`
	GameName     string `json:"game_name"`
	GameWebsite  string `json:"game_website"`
	LastUpdateId string `json:"last_update_id"`
}


func (g *GameData) getLastUpdate() UpdateData{
	var id string
	var url string

	switch g.GameId {

	case PUBG:
		rowData, _ := goquery.ParseUrl("http://www.playbattlegrounds.com/news.pu")
		temp := rowData.Find("#allList li a").Eq(0).Attr("href")
		url = "http://www.playbattlegrounds.com" + temp

		startIndex := strings.LastIndex(temp, "/") + 1
		endIndex := strings.LastIndex(temp, ".")
		tempId := []rune(temp)
		tempId = tempId[startIndex:endIndex]

		id = string(tempId)

	default:
		panic("i don't know this website!")
	}

	return UpdateData{id,url}
}


func (u *UpdateData) isUpdateHot(data GameData) bool {

	switch data.GameId {

	case PUBG:
		return u.Id  > data.LastUpdateId
	}


	return false
}


func pageGrabber() {

	hasUpdates := false

	for {

		for key, temp := range games {
			update := temp.getLastUpdate()

			if update.isUpdateHot(temp) {
				devide()
				log.Println("\tThere is new update for ", temp.GameName)

				temp.notifyUsersAboutUpdate(update)

				updatedGameData := games[key]
				updatedGameData.LastUpdateId = update.Id
				games[key]= updatedGameData

				hasUpdates = true
			}else{
				devide()
				log.Println("\tI haven't found updates for ", temp.GameName)
				devide()
			}
		}

		if hasUpdates {
			devide()
			log.Println("\tStart saving games")
			saveGamesData()
			log.Println("\tGames succesfully saved")
			devide()
			hasUpdates = false
		}

		time.Sleep(SLEEPING * time.Second)
	}

}