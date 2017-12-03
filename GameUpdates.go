package main

import (
	"log"
)

//Games id's
const(
	PUBG = 0
)


var games map[int]GameData


func main() {
	log.Println("\tStart")
	devide()

	log.Println("\tLoading Games Data")
	games = getGamesData()
	log.Println("\tSuccessfuly loaded GamesData(", len(games), ")")
	devide()

	log.Println("\tStarting preparing File System")
	prepareFileSystem(games)
	log.Println("\tSuccessfuly prepared File System")
	devide()

	log.Println("\tStarting Telegram Bot")
	go startBot()
	log.Println("\tSuccessfuly started Telegram Bot")
	devide()

	log.Println("\tStarting Page Grabber")
	pageGrabber()
}


func devide() {
	log.Println("----------------------------------")
}
