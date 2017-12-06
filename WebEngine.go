package main

import (
	"github.com/opesun/goquery"
	"strings"
	"log"
	"time"
)

const(
	SLEEPING = 60
)


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

	for {

		for key, temp := range games {
			update := temp.getLastUpdate()

			if update.isUpdateHot(temp) {
				devide()
				log.Println("\tThere is new update for", temp.GameShortName)

				go temp.notifyUsersAboutUpdate(update)

				updatedGameData := games[key]
				updatedGameData.LastUpdateId = update.Id
				updatedGameData.saveGamesData(db)
				games[key] = updatedGameData
			}else{
				devide()
				log.Println("\tI haven't found updates for ", temp.GameShortName)
				devide()
			}
		}

		time.Sleep(SLEEPING * time.Second)
	}

}