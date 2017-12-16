package db

import(
	_"github.com/jackc/pgx/stdlib"
	"database/sql"
	"log"
	"strconv"
	"encoding/json"
	"errors"
	"GameUpdates/data"
)


func getConnection() *sql.DB {
	if db == nil {
		db = openConnection()
	}

	return db
}


func openConnection() *sql.DB{
	db, err := sql.Open("pgx", connectionString)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Connection was opened")
	return db
}


func GetGamesData() map[int]*data.GameData {
	db := getConnection()

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

	result := make(map[int]*data.GameData)

	for rows.Next() {
		err = rows.Scan(&id, &shortName, &fullName, &lastUpdateId, &website)
		if err != nil{
			log.Fatal(err)
		}

		result[id] = &data.GameData{GameId: id, GameShortName: shortName, GameFullName : fullName, GameWebsite: website, LastUpdateId: lastUpdateId}
	}

	return result
}


func SaveGamesData(g *data.GameData) {
	db := getConnection()

	sqlQuery := "UPDATE games SET " + GAMES_SHORT_NAME + "='" + g.GameShortName + "',"
	sqlQuery += GAMES_FULL_NAME + "='" + g.GameFullName + "'," + GAMES_WEBSITE + "='" + g.GameWebsite + "'," + GAMES_LAST_UPDATE_ID + "='" + g.LastUpdateId + "' "
	sqlQuery += "WHERE " + GAMES_ID + "=" + strconv.Itoa(g.GameId)

	log.Println(sqlQuery)

	_, err := db.Exec(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
}


func GetAllUsers(g *data.GameData) []User{
	db := getConnection()
	sqlQuery := "SELECT * FROM users WHERE " + USERS_SUBSCRIBES + " like '[%" + strconv.Itoa(g.GameId) + "%]';"

	rows, err := db.Query(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		telegramId int64
		subscribes string
	)

	result := make([]User, 0, 10)

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


func GetUser(telegramId int64) (User, error) {
	db := getConnection()
	query := "select subscribes from users where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		var(
			rowSubscribes string
			subscribes    []int
		)

		err = rows.Scan(&rowSubscribes)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal([]byte(rowSubscribes), &subscribes)

		return User{TelegramId : telegramId, Subscribes:subscribes}, nil
	} else {
		return User{}, errors.New("no such user with id:" + strconv.FormatInt(telegramId, 10))
	}
}


func SubscribeUser(g *data.GameData, telegramId int64) bool{
	db := getConnection()

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


func UnSubscribeUser(g *data.GameData, telegramId int64) bool {
	db := getConnection()
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


func CloseConnection(){
	db.Close()
	log.Println("DB connection was closed")
}