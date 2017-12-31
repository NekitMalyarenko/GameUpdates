package web

import (
	"github.com/opesun/goquery"
	"strings"
	"data"
	"log"
	"telegram"
	"db"
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

	case data.RUST:
		rowData, _ := goquery.ParseUrl("https://rust.facepunch.com/blog/")
		root := rowData.Find(".yeargroup .column div .is-9")

		tempId := strings.TrimSpace(root.Find(".month").Eq(0).Text())
		endIndex := len(tempId) - 2
		id = tempId[:endIndex]

		url = "https://rust.facepunch.com" + rowData.Find(".is-10 div a").Attr("href")

	case data.CSGO:
		rowData, _ := goquery.ParseUrl("http://blog.counter-strike.net/index.php/category/updates/")
		temp := rowData.Find("#post_container .inner_post").Eq(0)

		id = temp.Find(".post_date").Text()
		id = strings.Replace(id, " ", "", -1)
		id = strings.Replace(id, "-", "", -1)

		url = temp.Find("h2 a").Attr("href")

	default:
		panic("i don't know this website!")
	}

	return UpdateData{id,url}
}


func (u *UpdateData) isUpdateHot(gameData *data.GameData) bool {

	switch gameData.GameId {

	case data.RUST:
		return u.Id != gameData.LastUpdateId

	default:
		return u.Id  > gameData.LastUpdateId
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
				db.GetDBManager().SaveGamesData(temp)
			}else{
				log.Println("\tI haven't found updates for", temp.GameShortName)
			}
		}

		time.Sleep( WEB_PARSER_SLEEPING * time.Second)
	}

}