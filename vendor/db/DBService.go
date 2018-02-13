package db

import(
	_"github.com/jackc/pgx/stdlib"
	"database/sql"
	"log"
	"strconv"
	"encoding/json"
	"data"
	"strings"
	"data/config"
)


func openConnection() *sql.DB{
	log.Println(configData.GetString(configData.CONNECTION_STRING))
	db, err := sql.Open("pgx", configData.GetString(configData.CONNECTION_STRING))
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Connection was opened")
	return db
}


func (manager *dbManager) GetGamesData() map[int]*data.GameData {
	sqlQuery := "SELECT * FROM games order by id asc;"

	rows, err := manager.db.Query(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id                 int
		shortName          string
		fullName           string
		website            string
		lastUpdateId       string
		countOfUsers       int
		countOfUniqueUsers int
	)

	result := make(map[int]*data.GameData)

	for rows.Next() {
		err = rows.Scan(&id, &shortName, &fullName, &lastUpdateId, &website, &countOfUsers, &countOfUniqueUsers)
		if err != nil{
			log.Fatal(err)
		}

		result[id] = &data.GameData{GameId: id, GameShortName: shortName, GameFullName : fullName, GameWebsite: website, LastUpdateId: lastUpdateId}
	}

	return result
}


func (manager *dbManager) SaveGamesData(g *data.GameData) {
	db := manager.db

	sqlQuery := "UPDATE games SET " + GAMES_SHORT_NAME + "='" + g.GameShortName + "',"
	sqlQuery += GAMES_FULL_NAME + "='" + g.GameFullName + "'," + GAMES_WEBSITE + "='" + g.GameWebsite + "'," + GAMES_LAST_UPDATE_ID + "='" + g.LastUpdateId + "' "
	sqlQuery += "WHERE " + GAMES_ID + "=" + strconv.Itoa(g.GameId)

	log.Println(sqlQuery)

	_, err := db.Exec(sqlQuery)
	if err != nil{
		log.Fatal(err)
	}
}


func (manager *dbManager) GetAllUsers(g *data.GameData) []User{
	db := manager.db
	var sqlQuery string

	if g != nil {
		sqlQuery = "SELECT * FROM users WHERE " + USERS_SUBSCRIBES + " like '[%" + strconv.Itoa(g.GameId) + "%]';"
	}else {
		sqlQuery = "SELECT * FROM users;"
	}

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


func (manager *dbManager) GetUser(telegramId int64) (User, error) {
	db := manager.db
	query := "select subscribes from users where " + USERS_TELEGRAM_ID + "=" + strconv.FormatInt(telegramId, 10)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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
		log.Println("no such user with id:" + strconv.FormatInt(telegramId, 10))
		return User{Subscribes: make([]int, 0)}, nil
	}
}


func (manager *dbManager) SubscribeUser(g *data.GameData, telegramId int64) bool{
	db := manager.db

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

		contains := func() bool {
			for _, temp := range subscribes{

				if temp == g.GameId {
					return true
				}
			}

			return false
		}

		if !contains() {
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


func (manager *dbManager) UnSubscribeUser(g *data.GameData, telegramId int64) bool {
	db := manager.db
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

			log.Println("value:", value)

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


func (manager *dbManager) SearchGame(request string) []int {
	db := manager.db

	request = strings.ToUpper(request)
	selectQuery := "select id from games where short_name ILIKE '%" + request + "%' or full_name ILIKE '%" + request + "%' order by id asc limit 4;"

	rows, err := db.Query(selectQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	result := make([]int, 0)
	var id int

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil{
			log.Fatal(err)
		}

		result = append(result, id)
	}

	return result
}


func (manager *dbManager) UpdateGamesSubscribersData() {
	db := manager.db
	//var temp int

	for _, val := range data.GetGames() {
		selectUsersNumberQuery := "select from users count(id) where subscribes like '[%" + strconv.Itoa(val.GameId) + "%]';"
		selectUniqueUsersNumberQuery := "select from users count(id) where subscribes like '[" + strconv.Itoa(val.GameId) + "]';"

		rows, err := db.Exec(selectUsersNumberQuery)
		if err != nil {
			log.Fatal(err)
		}
		number, _ := rows.RowsAffected()
		updateUsersNumber := "update games set number_of_users = " + strconv.FormatInt(number, 10) + " where id = " + strconv.Itoa(val.GameId) + ";"
		db.Exec(updateUsersNumber)

		rows, err = db.Exec(selectUniqueUsersNumberQuery)
		if err != nil {
			log.Fatal(err)
		}
		number, _ = rows.RowsAffected()
		updateUsersNumber = "update games set number_of_unique_users = " + strconv.FormatInt(number, 10) + " where id = " + strconv.Itoa(val.GameId) + ";"
		db.Exec(updateUsersNumber)
	}


}


func CloseConnection(){
	manager.db.Close()
	log.Println("DB connection was closed")
}