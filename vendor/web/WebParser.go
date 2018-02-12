package web

import (
	"github.com/opesun/goquery"
	"strings"
	"data"
	"log"
	"telegram"
	"db"
	"time"
	"data/config"
	"fmt"
)

type UpdateData struct {
	Id  string
	Url string
}


func getLastUpdate(g *data.GameData)  (UpdateData, bool){
	id := ""
	url := ""
	isWebsiteDown := false

	switch g.GameId {

	case data.PUBG:
		rowData, _ := goquery.ParseUrl("http://www.playbattlegrounds.com/news.pu")

		if rowData != nil && len(rowData.Text()) != 0 {
			temp := rowData.Find("#allList li a").Eq(0).Attr("href")
			url = "http://www.playbattlegrounds.com" + temp

			startIndex := strings.LastIndex(temp, "/") + 1
			endIndex := strings.LastIndex(temp, ".")
			tempId := []rune(temp)
			tempId = tempId[startIndex:endIndex]

			id = string(tempId)
		} else {
			isWebsiteDown = true
		}

	case data.GTA:
		rowData, _ := goquery.ParseUrl("http://steamcommunity.com/games/271590/announcements/")
		if rowData != nil && len(rowData.Text()) != 0 {
			url = rowData.Find("#announcementsContainer .announcement .large_title").Eq(0).Attr("href")

			runes := []byte(url)
			startIndex, endIndex := strings.LastIndex(url, "/") + 1, len(url)

			id = string(runes[startIndex:endIndex])
		} else {
			isWebsiteDown = true
		}

	case data.RUST:
		rowData, _ := goquery.ParseUrl("https://rust.facepunch.com/blog/")
		if rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find(".yeargroup .column div .is-9")

			tempId := strings.TrimSpace(root.Find(".month").Eq(0).Text())
			endIndex := len(tempId) - 2
			id = tempId[:endIndex]

			url = "https://rust.facepunch.com" + rowData.Find(".is-10 div a").Attr("href")
		} else {
			isWebsiteDown = true
		}

	case data.CSGO:
		rowData, _ := goquery.ParseUrl("http://blog.counter-strike.net/index.php/category/updates/")
		if rowData != nil && len(rowData.Text()) != 0 {
			temp := rowData.Find("#post_container .inner_post").Eq(0)

			id = temp.Find(".post_date").Text()
			id = strings.Replace(id, " ", "", -1)
			id = strings.Replace(id, "-", "", -1)

			url = temp.Find("h2 a").Attr("href")
		} else {
			isWebsiteDown = true
		}

	case data.OVERWATCH:
		rowData, _ := goquery.ParseUrl("https://playoverwatch.com/en-us/game/patch-notes/pc")
		if rowData != nil && rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find(".patch-notes-default div .lg-9 .patch-notes-body").First()

			id = root.Attr("id")
			startIndex := strings.Index(id, "-") + 1
			id = id[startIndex:]

			url = "https://playoverwatch.com/en-us/game/patch-notes/pc"
		} else {
			isWebsiteDown = true
		}
		break

	case data.GOVNO:
		rowData, _ := goquery.ParseUrl("http://www.dota2.com/news/updates/")
		if rowData != nil && len(rowData.Text()) != 0{
			root := rowData.Find("#mainLoop div")

			id = root.Attr("id")
			startIndex := strings.Index(id, "-") + 1
			id = id[startIndex:]

			url = root.Find("h2 a").Attr("href")
		} else {
			isWebsiteDown = true
		}
		break

	case data.SQUAD:
		rowData, _ := goquery.ParseUrl("http://joinsquad.com/")

		if rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find("#updates .updates-content-box .update").Eq(0)

			id = root.Find("a").Attr("href")
			url = "http://joinsquad.com" + id

			startIndex := strings.LastIndex(id, "=") + 1
			id = id[startIndex:]
		} else {
			isWebsiteDown = true
		}
		break

	case data.RAINBOW:
		rowData, _ := goquery.ParseUrl("http://store.steampowered.com/news/?appids=359550")
		if rowData != nil && len(rowData.Text()) != 0 {
			rowData = rowData.Find("#news div").Eq(1)
			url = rowData.Find("div .headline .posttitle a").Attr("href")

			startIndex := strings.LastIndex(url, "/") + 1
			id = string(url[startIndex:])
		} else {
			isWebsiteDown = true
		}
		break

	case data.TEAMFORTESS:
		rowData, _ := goquery.ParseUrl("http://www.teamfortress.com/?tab=updates")
		if rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find("#leftColPosts a").Eq(0)

			id = root.Attr("href")
			startIndex := strings.Index(id, "=") + 1
			id = id[startIndex:]

			url = "http://www.teamfortress.com/" + root.Attr("href")
		} else {
			isWebsiteDown = true
		}
		break

	case data.LOL:
		rowData, _ := goquery.ParseUrl("https://playhearthstone.com/en-us/blog/")
		if rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find("#blog-articles li").Eq(1)

			id = root.Attr("data-id")
			url = "https://playhearthstone.com" + root.Find(".media__bd .article-title a").Attr("href")
		} else {
			isWebsiteDown = true
		}
		break

	case data.FORZA_MOTOSPORT:
		rowData, _ := goquery.ParseUrl("https://forzamotorsport.net/en-US/news")
		if rowData != nil && len(rowData.Text()) != 0 {
			root := rowData.Find(".news_list .common_item").Eq(0)

			url = "https://forzamotorsport.net" + root.Find(".media_container .media_element a").Attr("href")
			id = root.Find(".author_date .date").Text()
		} else {
			isWebsiteDown = true
		}
		break

	default:
		panic("i don't know this website!")
	}

	return UpdateData{id,url}, isWebsiteDown
}


func (u *UpdateData) isUpdateHot(gameData *data.GameData) bool {

	switch gameData.GameId {

	case data.RUST, data.RAINBOW:
		return u.Id != gameData.LastUpdateId

	default:
		return u.Id  > gameData.LastUpdateId
	}

	return false
}


func PageGrabber() {
	games := data.GetGames()
	sleeping := time.Duration(configData.GetInt(configData.WEB_PARSER_TIMEOUT))

	for {
		for _, temp := range games {
			update, isWebsiteDown := getLastUpdate(temp)

			if !isWebsiteDown {

				if update.isUpdateHot(temp) {
					log.Println("\tThere is new update for", temp.GameShortName)
					telegram.SendMessageToMe("There is new update for " + temp.GameShortName + "\t" + update.Id + "\t" + update.Url)

					go telegram.NotifyUsersAboutUpdate(temp, update.Url)

					temp.LastUpdateId = update.Id
					db.GetDBManager().SaveGamesData(temp)
				} else {
					fmt.Println("\tI haven't found updates for", temp.GameShortName)
				}
			} else {
				log.Println("\t" + temp.GameShortName, "is down")
				//telegram.SendMessageToMe(temp.GameShortName + " is down")
			}
		}

		time.Sleep(sleeping * time.Second)
	}

}
