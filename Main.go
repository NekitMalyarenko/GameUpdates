package main

import (
	"log"
	"db"
	"time"
)


func main() {

	test()

	/*log.Println("Start")

	log.Println("Getting games")
	data.SetGames(db.GetDBManager().GetGamesData())

	defer db.CloseConnection()
	log.Println("Successfuly got games")

	log.Println("Starting Telegram Bot")
	go telegram.StartBot()
	log.Println("Successfuly started Telegram Bot")

	time.Sleep(5 * time.Second)

	log.Println("Starting Page Grabber")
	web.PageGrabber()*/
}


func test() {
	log.Println("Test 1")
	db.GetDBManager().GetGamesData()
	go db.GetDBManager().GetUser(0)
	go db.GetDBManager().GetGamesData()
	go db.GetDBManager().GetGamesData()

	time.Sleep(120 * time.Second)

	db.CloseConnection()
}

