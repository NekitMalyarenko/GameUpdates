package main

import (
	"log"
	"database/sql"
	"time"
)

type GameData struct {
	GameId        int    `json:"game_id"`
	GameShortName string `json:"game_name"`
	GameFullName  string `json:"game_name"`
	GameWebsite   string `json:"game_website"`
	LastUpdateId  string `json:"last_update_id"`
}


type UpdateData struct {
	Id  string
	Url string
}


//Games id's
const(
	PUBG = 1
)

var games map[int]GameData

var db *sql.DB


func main() {

	log.Println("\tStart")
	devide()

	db = openConnection()
	defer closeConnection(db)

	log.Println("\tLoading Games Data")
	games = getGamesData(db)
	log.Println("\tSuccessfuly loaded GamesData(", games, ")")
	devide()

	log.Println("\tStarting Telegram Bot")
	go startBot()
	log.Println("\tSuccessfuly started Telegram Bot")
	devide()

	time.Sleep(5 * time.Second)

	log.Println("\tStarting Page Grabber")
	pageGrabber()
}


func devide() {
	log.Println("----------------------------------")
}

func CheckError(err error){
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
