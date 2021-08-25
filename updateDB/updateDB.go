package updateDB

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type ChannelWrap struct {
	ChannelID string
	TeamToken string
}
type Board struct {
	ID           int
	Name         string
	NameKor      string
	Link         string
	IsCsBoard    bool
	LastNotified time.Time
}

func goDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading .env file")
	}
	return os.Getenv(key)
}
func openDB() *sql.DB {
	db, err := sql.Open("mysql", goDotEnvVar("MYSQL_ID")+":"+goDotEnvVar("MYSQL_PW")+"@tcp("+goDotEnvVar("MYSQL_HOST")+")/"+goDotEnvVar("MYSQL_DB")+"?parseTime=true")
	ErrCheck(err)
	return db
}
func ErrCheck(e error) bool {
	if e != nil {
		panic(e)
		return true
	}
	return false
}
func GetBoardByName(boardName string) []Board {
	db := openDB()
	rows, err := db.Query("SELECT * FROM board_list WHERE name=?", boardName)
	ErrCheck(err)
	defer db.Close()
	defer rows.Close()

	ret := []Board{}
	item := Board{}
	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Name, &item.NameKor, &item.Link, &item.IsCsBoard, &item.LastNotified)
		ErrCheck(err)
		ret = append(ret, item)
	}

	return ret
}
func FindBoards() []Board {
	db := openDB()
	rows, err := db.Query("SELECT * FROM board_list")
	ErrCheck(err)
	defer db.Close()
	defer rows.Close()

	ret := []Board{}
	item := Board{}
	for rows.Next() {
		err = rows.Scan(&item.ID, &item.Name, &item.NameKor, &item.Link, &item.IsCsBoard, &item.LastNotified)
		ErrCheck(err)
		ret = append(ret, item)
	}

	return ret
}
func InsertChannel(teamToken string, teamID string, channelID string, boardType string) bool {
	if boardType == "" {
		return false
	}
	db := openDB()
	rows, err := db.Query("SELECT EXISTS (SELECT * FROM subscription WHERE channelID=? AND boardType=(SELECT id FROM board_list WHERE name=?)) AS chk", channelID, boardType)
	ErrCheck(err)
	defer db.Close()
	defer rows.Close()

	var ret bool
	for rows.Next() {
		err = rows.Scan(&ret)
		ErrCheck(err)
	}
	if ret == false {
		res, err := db.Query("SELECT EXISTS (SELECT * FROM teamInfo WHERE teamID = ?) AS chk", teamID)
		ErrCheck(err)
		defer res.Close()

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
		defer res.Close()

		var idx int
		for res.Next() {
			err = res.Scan(&idx)
			ErrCheck(err)
		}

		res, err = db.Query("INSERT INTO subscription(boardType, channelID, teamID) VALUES(?, ?, ?)", boardType, channelID, idx)
		ErrCheck(err)
		defer res.Close()

		return true
	}
	return false
}
func DeleteChannel(channelID string, boardType string) bool {
	db := openDB()
	rows, err := db.Query("SELECT EXISTS (SELECT * FROM subscription WHERE channelID=? AND boardType=(SELECT id FROM board_list WHERE name=?)) AS chk", channelID, boardType)
	ErrCheck(err)
	defer db.Close()
	defer rows.Close()

	var ret bool
	for rows.Next() {
		err = rows.Scan(&ret)
		ErrCheck(err)
	}
	if ret == true {
		res, err := db.Query("DELETE FROM subscription WHERE channelID=? AND boardType=(SELECT id FROM board_list WHERE name=?)", channelID, boardType)
		ErrCheck(err)
		defer res.Close()
		return true
	}
	return false
}

func FindChannels(boardType []string) []ChannelWrap {
	var channelIDList []ChannelWrap
	db := openDB()
	defer db.Close()

	for _, board := range boardType {
		rows, err := db.Query("SELECT channelID, teamInfo.teamToken FROM subscription JOIN teamInfo ON subscription.teamID = teamInfo.idx WHERE boardType=(SELECT id FROM board_list WHERE name=?)", board)
		ErrCheck(err)
		defer rows.Close()

		for rows.Next() {
			tmp := new(ChannelWrap)
			err = rows.Scan(&tmp.ChannelID, &tmp.TeamToken)
			ErrCheck(err)
			channelIDList = append(channelIDList, *tmp)
		}
	}

	return channelIDList
}
func SetTeamToken(teamID string, teamToken string) bool {
	db := openDB()

	rows, err := db.Query("SELECT EXISTS (SELECT * FROM teamInfo WHERE teamID = ?) AS chk", teamID)
	ErrCheck(err)
	defer rows.Close()
	defer db.Close()

	var chk bool
	for rows.Next() {
		err = rows.Scan(&chk)
		ErrCheck(err)
	}

	// TeamID is already in DB
	if chk == true {
		rows, err = db.Query("UPDATE teamInfo SET teamToken = ? WHERE teamID = ?", teamToken, teamID)
		ErrCheck(err)
		defer db.Close()
		rows.Close()
	} else {
		// Newly added TeamID
		rows, err = db.Query("INSERT INTO teamInfo(teamID, teamToken) VALUES(?, ?)", teamID, teamToken)
		ErrCheck(err)
		defer db.Close()
		rows.Close()
	}
	return true
}
func GetTeamToken(teamID string) string {
	db := openDB()

	rows, err := db.Query("SELECT teamToken FROM teamInfo WHERE teamID = ?", teamID)
	ErrCheck(err)
	defer db.Close()
	defer rows.Close()

	var Token string
	for rows.Next() {
		err = rows.Scan(&Token)
		ErrCheck(err)
	}
	return Token
}
