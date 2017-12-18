package web

import (
	"github.com/opesun/goquery"
	"strings"
	"github.com/NekitMalyarenko/GameUpdates/data"
	"log"
	"github.com/NekitMalyarenko/GameUpdates/telegram"
	"github.com/NekitMalyarenko/GameUpdates/db"
	"time"
)


type UpdateData struct {
	Id  string
	Url string
}

const(
	WEB_PARSER_SLEEPING = 60
)


func getLastUpdate(g *data.GameData) UpdateData{
	var id string
	var url string

	switch g.GameId {

	case data.PUBG:
		rowData, _ := goquery.ParseUrl("http://www.playbattlegrounds.com/news.pu")
		temp := rowData.Find("#allList li a").Eq(0).Attr("href")
		url = "http://www.playbattlegrounds.com" + temp

		startIndex := strings.LastIndex(temp, "/") + 1
		endIndex := strings.LastIndex(temp, ".")
		tempId := []rune(temp)
		tempId = tempId[startIndex:endIndex]

		id = string(tempId)

	case data.GTA:
		rowData, _ := goquery.ParseUrl("http://steamcommunity.com/games/271590/announcements/")
		url := rowData.Find("#announcementsContainer .announcement .large_title").Eq(0).Attr("href")

		runes := []byte(url)
		startIndex, endIndex := strings.LastIndex(url, "/") + 1, len(url)

		id = string(runes[startIndex:endIndex])

	default:
		panic("i don't know this website!")
	}

	return UpdateData{id,url}
}


func (u *UpdateData) isUpdateHot(data *data.GameData) bool {

	switch data.GameId {
		default:
			return u.Id  > data.LastUpdateId
	}

	return false
}


func PageGrabber() {

	games := data.GetGames()

	for {
		for _, temp := range games {
			update := getLastUpdate(temp)

			if update.isUpdateHot(temp) {
				log.Println("\tThere is new update for", temp.GameShortName)

				go telegram.NotifyUsersAboutUpdate(temp, update.Url)

				temp.LastUpdateId = update.Id
				db.SaveGamesData(temp)
			}else{
				log.Println("\tI haven't found updates for", temp.GameShortName)
			}
		}

		time.Sleep( WEB_PARSER_SLEEPING * time.Second)
	}

}