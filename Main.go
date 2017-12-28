package main

import (
	"log"
	"db"
	"time"
	"data"
	"telegram"
	"web"
)


func main() {
	log.Println("Start")

	log.Println("Getting games")
	data.SetGames(db.GetDBManager().GetGamesData())

	defer db.CloseConnection()
	log.Println("Successfuly got games")

	log.Println("Starting Telegram Bot")
	go telegram.StartBot()
	log.Println("Successfuly started Telegram Bot")

	time.Sleep(5 * time.Second)

	log.Println("Starting Page Grabber")
	web.PageGrabber()
}


