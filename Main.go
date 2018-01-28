package main

import (
	"os"
	"log"
	"data"
	"db"
	"telegram"
	"time"
	"web"
)


func main() {
	log.Println("\tSTART")

	log.Println("Getting games")
	data.SetGames(db.GetDBManager().GetGamesData())
	defer db.CloseConnection()

	log.Println("Starting Telegram Bot")
	go telegram.StartBot()
	log.Println("Successfuly started Telegram Bot")

	time.Sleep(5 * time.Second)

	log.Println("Starting Page Grabber")
	web.PageGrabber()
}