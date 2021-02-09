package updateDB

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ChannelWrap struct {
	ChannelID string
	TeamToken string
}

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
func AddChannel(teamToken string, teamID string, channelID string, boardType string) bool {
	if boardType == "" {
		return false
	}

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
		res, err := db.Query("SELECT EXISTS (SELECT * FROM teamInfo WHERE teamID = ?) AS chk", teamID)
		ErrCheck(err)
		var chk bool
		for res.Next() {
			err = res.Scan(&chk)
			ErrCheck(err)
		}
		if !chk {
			res, err := db.Query("INSERT INTO teamInfo(teamID, teamToken) VALUES(?, ?)", teamID, teamToken)
			ErrCheck(err)
			res.Close()
		}
		res, err = db.Query("SELECT idx FROM teamInfo WHERE teamID = ?", teamID)
		ErrCheck(err)
		var idx int
		for res.Next() {
			err = res.Scan(&idx)
			ErrCheck(err)
		}
		// fmt.Println(idx)
		res, err = db.Query("INSERT INTO "+boardType+"Notice(channelID, teamIdx) VALUES(?, ?)", channelID, idx)
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
func GetChannels(boardType string) []ChannelWrap {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)
	rows, err := db.Query("SELECT channelID, teamInfo.teamToken FROM " + boardType + "Notice JOIN teamInfo ON " + boardType + "Notice.teamIdx = teamInfo.idx")
	ErrCheck(err)
	defer db.Close()

	var channelIDList []ChannelWrap
	for rows.Next() {
		tmp := new(ChannelWrap)
		err = rows.Scan(&tmp.ChannelID, &tmp.TeamToken)
		ErrCheck(err)
		channelIDList = append(channelIDList, *tmp)
	}
	return channelIDList
}
func SetTeamToken(teamID string, teamToken string) bool {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)

	rows, err := db.Query("SELECT EXISTS (SELECT * FROM teamInfo WHERE teamID = ?) AS chk", teamID)
	ErrCheck(err)

	var chk bool
	for rows.Next() {
		err = rows.Scan(&chk)
		ErrCheck(err)
	}

	// TeamID is alreay in DB
	if chk == true {
		return false
	}
	// Newly added TeamID
	rows, err = db.Query("INSERT INTO teamInfo(teamID, teamToken) VALUES(?, ?)", teamID, teamToken)
	ErrCheck(err)
	defer db.Close()
	rows.Close()

	return true
}
func GetTeamToken(teamID string) string {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp(leed.at:3306)/"+goDotEnvVar("MYSQL_DB"))
	ErrCheck(err)
	rows, err := db.Query("SELECT teamToken FROM teamInfo WHERE teamID = ?", teamID)
	ErrCheck(err)
	defer db.Close()

	var Token string
	for rows.Next() {
		err = rows.Scan(&Token)
		ErrCheck(err)
	}
	return Token
}
