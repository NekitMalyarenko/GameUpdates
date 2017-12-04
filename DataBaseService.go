package main

import(
	_"github.com/lib/pq"
	"database/sql"
	"os"
	"log"
)


func test() {
	db, err := sql.Open("postgres", os.Getenv("host=ec2-184-72-247-126.compute-1.amazonaws.com port=5432 user=akbrukfjbdaijj password=04bfa834c09cd695ddfaa3dae1ff5f2122cfae835c7eac9654447693331a6ceb dbname=ddhcg53rvj0vis sslmode=disable"))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE users (telegramId INTEGER PRIMARY KEY,subscribes SMALLINT []);")
	if err != nil {
		log.Fatal(err)
	}
}


