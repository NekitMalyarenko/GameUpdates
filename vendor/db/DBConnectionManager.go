package db

import (
	"database/sql"
	"sync"
	"fmt"
)

type User struct {
	TelegramId int64
	Subscribes []int
}

type dbManager struct {
	db *sql.DB
}


const(
	GAMES_ID              = "id"
	GAMES_SHORT_NAME      = "short_name"
	GAMES_FULL_NAME       = "full_name"
	GAMES_WEBSITE         = "game_website"
	GAMES_LAST_UPDATE_ID  = "last_update_id"

	USERS_TELEGRAM_ID     = "telegram_id"
	USERS_SUBSCRIBES      = "subscribes"
)

var (
	manager *dbManager
	once sync.Once
)


func GetDBManager() *dbManager {
	once.Do(func() {
		manager = &dbManager{
			db:openConnection(),
		}
		fmt.Println("New Connection")
	})

	return manager
}