package updateDB

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func goDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading .env file")
	}
	return os.Getenv(key)
}
func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}
func AddChannel(channelID string, boardType string) bool {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)
	rows, err := db.Query("SELECT EXISTS (SELECT * FROM "+boardType+"Notice WHERE channelID = ?) AS chk", channelID)
	ErrCheck(err)
	defer db.Close()

	var ret bool
	for rows.Next() {
		err = rows.Scan(&ret)
		ErrCheck(err)
	}
	if ret == false {
		res, err := db.Query("INSERT INTO "+boardType+"Notice(channelID) VALUES(?)", channelID)
		ErrCheck(err)
		res.Close()
		return true
	}
	return false
}
func RemoveChannel(channelID string, boardType string) bool {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)
	rows, err := db.Query("SELECT EXISTS (SELECT * FROM "+boardType+"Notice WHERE channelID = ?) AS chk", channelID)
	ErrCheck(err)
	defer db.Close()

	var ret bool
	for rows.Next() {
		err = rows.Scan(&ret)
		ErrCheck(err)
	}
	if ret == true {
		res, err := db.Query("DELETE FROM "+boardType+"Notice WHERE channelID = ?", channelID)
		ErrCheck(err)
		res.Close()
		return true
	}
	return false
}
func GetChannels(boardType string) []string {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)
	rows, err := db.Query("SELECT channelID FROM " + boardType + "Notice")
	ErrCheck(err)
	defer db.Close()

	var id string
	var channelIDList []string
	for rows.Next() {
		err = rows.Scan(&id)
		ErrCheck(err)
		channelIDList = append(channelIDList, id)
	}
	return channelIDList
}
