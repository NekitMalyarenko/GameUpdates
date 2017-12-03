package main

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"errors"
	"log"
	"strconv"
)


const (
	GAMESDATAPATH = "games_data.json"
	DATAPATH      = "data/"
)


func prepareFileSystem(data map[int]GameData ) bool{
	os.Mkdir(DATAPATH, os.ModePerm)

	var path string
	result := true

	for _, game := range data {

		path = DATAPATH + game.GameName + ".json"

		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			log.Println("Creating :" + path)
			_, err := os.Create(path)

			if err != nil {
				panic(err)
			}else{
				ioutil.WriteFile(path, []byte("[]"), os.ModeAppend)
			}

			result = false
		}
	}

	return result
}


func getGamesData() map[int]GameData {
	result := make(map[int]GameData)

	data := make([]GameData, 0, 1)
	log.Println("GamesDataPath:", GAMESDATAPATH)
	raw, err := ioutil.ReadFile(GAMESDATAPATH)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &data)

	for _, temp := range data {
		result[temp.GameId] = temp
	}

	return result
}


func getPath(gameId int) (string, error) {

	for _, temp := range games{
		if temp.GameId == gameId{
			return DATAPATH + temp.GameName + ".json", nil
		}
	}

	return "", errors.New("undefined gameId")
}


func getUsers(gameId int) []string{
	result := make([]string, 1, 2)
	path, err := getPath(gameId)

	if err == nil {
		raw, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		json.Unmarshal(raw, &result)
		return result

	}

	panic(err)
}


func addUser(gameId int, userId int) {
	str := strconv.Itoa(userId)

	data := getUsers(gameId)
	data = append(data, str)

	path, err := getPath(gameId)

	if err != nil {
		panic(err)
	}

	jsonData, _ := json.Marshal(data)

	error := ioutil.WriteFile(path, jsonData, os.ModeAppend)

	if error != nil {
		panic(error)
	}
}


func removeUser(gameId int, userId int){
	str := strconv.Itoa(userId)

	data := getUsers(gameId)

	for index, temp := range data {
		if temp == str {
			data = append(data[:index], data[index+1:]...)
			break
		}
	}

	path, err := getPath(gameId)

	if err != nil {
		panic(err)
	}

	jsonData, _ := json.Marshal(data)

	error := ioutil.WriteFile(path, jsonData, os.ModeAppend)

	if error != nil {
		panic(error)
	}
}


func saveGamesData(){
	result := make([]GameData, 0, 1)

	for _, temp := range games {
		result = append(result, temp)
	}

	data, err := json.Marshal(result)
	if err != nil{
		panic(err)
	}

	err = ioutil.WriteFile(GAMESDATAPATH, data, os.ModeAppend)
	if err != nil {
		panic(err)
	}
}