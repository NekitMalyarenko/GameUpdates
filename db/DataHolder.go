package db

import "database/sql"

type User struct {
	TelegramId int64
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

var (
	connectionString = "user=knpamodhrsrykd password=f426870c9669c69c1b5db08f02e2b268851a4d31a194a59b9bd2bf96ac3bd28f host=ec2-54-235-76-111.compute-1.amazonaws.com port=5432 database=dcnrf1jkrmd6k7 sslmode=require"
	db *sql.DB
)