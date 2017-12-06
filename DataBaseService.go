package main

import(
	_"github.com/jackc/pgx/stdlib"
	"database/sql"
	"log"
	"strconv"
	"encoding/json"
)

type User struct {
	TelegramId int
	Subscribes []int
}

const(
	GAMES_ID =             "id"
	GAMES_SHORT_NAME =     "short_name"
	GAMES_FULL_NAME =      "full_name"
	GAMES_WEBSITE =        "game_website"
	GAMES_LAST_UPDATE_ID = "last_update_id"

	USERS_TELEGRAM_ID =    "telegram_id"
	USERS_SUBSCRIBES =     "subscribes"
)

var connectionString = "user=knpamodhrsrykd password=f426870c9669c69c1b5db08f02e2b268851a4d31a194a59b9bd2bf96ac3bd28f host=ec2-54-235-76-111.compute-1.amazonaws.com port=5432 database=dcnrf1jkrmd6k7 sslmode=require"


func openConnection() *sql.DB{
	db, err := sql.Open("pgx", connectionString)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("\tConnection was opened")
	return db
}


func getGamesData(db *sql.DB) map[int]GameData {
	sqlQuery := "SELECT * FROM games;"

	rows, err := db.Query(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id           int
		shortName    string
		fullName     string
		website      string
		lastUpdateId string
	)

	result := make(map[int]GameData)

	for rows.Next() {
		err = rows.Scan(&id, &shortName, &fullName, &lastUpdateId, &website)
		if err != nil{
			log.Fatal(err)
		}

		result[id] = GameData{GameId: id, GameShortName: shortName, GameFullName : fullName, GameWebsite: website, LastUpdateId: lastUpdateId}
	}

	return result
}


func (g *GameData) saveGamesData(db *sql.DB) {
	sqlQuery := "UPDATE games SET " + GAMES_SHORT_NAME + "='" + g.GameShortName + "',"
	sqlQuery += GAMES_FULL_NAME + "='" + g.GameFullName + "'," + GAMES_WEBSITE + "='" + g.GameWebsite + "'," + GAMES_LAST_UPDATE_ID + "='" + g.LastUpdateId + "' "
	sqlQuery += "WHERE " + GAMES_ID + "=" + strconv.Itoa(g.GameId)

	log.Println(sqlQuery)

	_, err := db.Exec(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
}


func (g *GameData) getAllUsers(db *sql.DB) []User{
	sqlQuery := "SELECT * FROM users WHERE " + USERS_SUBSCRIBES + " like '[%" + strconv.Itoa(g.GameId) + "%]';"
	log.Println(sqlQuery)
	rows, err := db.Query(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		telegramId int
		subscribes string
	)

	result := make([]User, 1, 10)

	for i := 0;rows.Next();i++ {
		err = rows.Scan(&telegramId, &subscribes)
		if err != nil {
			log.Fatal(err)
		}

		temp := make([]int, 1)

		err = json.Unmarshal([]byte(subscribes), &temp)
		if err != nil{
			log.Fatal(err)
		}

		result = append(result, User{TelegramId : telegramId, Subscribes : temp})
	}

	return result
}


func (g *GameData) subscribeUser(db *sql.DB, telegramId int64) bool{
	selectQuery := "select subscribes from users where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10) + ";"
	insertQuery := "insert into users(" + USERS_TELEGRAM_ID + "," + USERS_SUBSCRIBES + ")" + " values('" + strconv.FormatInt(telegramId, 10) + "','[" + strconv.Itoa(g.GameId) + "]');"

	rows, err := db.Query(selectQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var subscribes []int

	if rows.Next(){
		var rowSubscribes string
		err = rows.Scan(&rowSubscribes)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(rowSubscribes), &subscribes)

		contains := false
		for _, temp := range subscribes{

			if temp == g.GameId {
				contains = true
			}
		}

		if !contains {
			subscribes = append(subscribes, g.GameId)
			resSub, _ := json.Marshal(subscribes)

			updateQuery := "update users set " + USERS_SUBSCRIBES + "='" + string(resSub) + "' where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10)
			_, err = db.Exec(updateQuery)

			if err == nil {
				return true
			} else {
				log.Fatal(err)
				return false
			}
		} else {
			return false
		}

	}else{
		_, err = db.Exec(insertQuery)
		if err == nil{
			return true
		} else {
			log.Fatal(err)
			return false
		}
	}

	return false
}


func (g *GameData) unSubscribeUser(db *sql.DB, telegramId int64) bool {
	selectQuery := "select subscribes from users where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10) + ";"

	rows, err := db.Query(selectQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var subscribes []int
		var rowSubscribes string

		err = rows.Scan(&rowSubscribes)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(rowSubscribes), &subscribes)
		temp := -1

		for index, value := range subscribes  {

			if value == g.GameId {
				temp = index
				break
			}

		}

		if temp != -1 {
			subscribes = append(subscribes[:temp], subscribes[temp+1:]...)
			resSub, _ := json.Marshal(subscribes)

			updateQuery := "update users set " + USERS_SUBSCRIBES + "='" + string(resSub) + "' where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10) + ";"
			_, err = db.Exec(updateQuery)
			if err != nil{
				log.Fatal(err)
			}

			return true
		} else {
			return false
		}

	} else {
		return false
	}
}


func closeConnection(db *sql.DB){
	db.Close()
	log.Println("\tDB connection was closed")
}