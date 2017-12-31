package data

type GameData struct {
	GameId        int    `json:"game_id"`
	GameShortName string `json:"game_name"`
	GameFullName  string `json:"game_name"`
	GameWebsite   string `json:"game_website"`
	LastUpdateId  string `json:"last_update_id"`
}


const(
	PUBG  = 1
	GTA   = 2
	RUST  = 3
	CSGO  = 4
)


var games map[int]*GameData = nil


func SetGames(g map[int]*GameData) {
	games = g
}


func GetGames() map[int]*GameData {
	return games
}


func UpdateGame(id int, lastUpdateId string) {
	(games[id]).LastUpdateId = lastUpdateId
}


func GetGame(id int) *GameData {
	return games[id]
}