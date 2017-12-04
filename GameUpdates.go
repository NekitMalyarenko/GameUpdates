package main

import (
	"log"
	"io/ioutil"
	"os"
	"strings"
)

//Games id's
const(
	PUBG = 0
)

var games map[int]GameData

//var root = "D:/Projects/go/GameUpdates/src"


func main() {
	log.Println("\tStart")
	test()
	log.Println("\tEND")

	/*log.Println("\tStart")
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
	pageGrabber()*/
}


func devide() {
	log.Println("----------------------------------")
}


func checkFile(name string) {

	log.Println(name)

	file, err := os.Open(name)
	checkError(err)
	defer file.Close()

	fileInfo, err := file.Stat()

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(name)
		checkError(err)

		for _, temp := range files {

			if strings.Contains(temp.Name(), "."){
				if strings.Contains(temp.Name(), ".go"){
					checkFile(name + "/" + temp.Name())
				}
			}else {
				checkFile(name + "/" + temp.Name())
			}
		}
	} else {
		row, err := ioutil.ReadFile(name)
		checkError(err)

		data := string(row[:])
		startIndex := strings.Index(data, "import (") + 8
		endIndex := strings.Index(data, ")")

		if startIndex != -1 && endIndex != -1 {

			temp := []rune(data)
			data = string(temp[startIndex:endIndex])

			if strings.Contains(data, "\"\"") {
				log.Println("-----------------------------------")
				log.Println("-------------NOT CLEAR-------------")
				log.Println("-----------------------------------")
			} else {
				log.Println("CLEAR")
			}
		} else {
			log.Println("-------------NO IMPORTS-------------")
		}
	}

}


func checkError(err error){
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
