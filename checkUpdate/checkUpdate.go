package checkupdate

import (
	"github.com/anaskhan96/soup"
	"log"
	SlackApi "slackApi"
	"time"
	"updateDB"
)

type BoardItems struct {
	Title string
	Url   string
}
type PKID struct {
	Strong struct {
		Main    string `json:"Main"`
		Underg  string `json:"Underg"`
		Grad    string `json:"Grad"`
		General string `json:"General"`
		Job     string `json:"Job"`
		Sgcs    string `json:"Sgcs"`
	} `json:"strong"`
	General struct {
		Main    string `json:"Main"`
		Underg  string `json:"Underg"`
		Grad    string `json:"Grad"`
		General string `json:"General"`
		Job     string `json:"Job"`
		Sgcs    string `json:"Sgcs"`
	} `json:"general"`
}
type ChannelWrap struct {
	ChannelID string
	TeamToken string
}

func ErrCheck(e error) {
	if e != nil {
		log.Panic(e)
	}
}
func FetchIntervally() {
	t := time.NewTicker(30 * time.Minute)
	// run intervally
	for range t.C {
		boards := updateDB.FindBoards()
		CheckBoardsAndNotify(boards)
	}
}
func CheckBoardsAndNotify(boards []updateDB.Board) {
	for _, board := range boards {
		var channelIDList = updateDB.FindChannels([]string{board.Name})
		items := FindUpdatedArticle(board.Link, board.IsCsBoard, board.LastNotified)
		for _, receiverInfo := range channelIDList {
			for _, val := range items {
				SlackApi.SendMessage(receiverInfo, board.NameKor, val.Title, val.Url)
			}
		}
	}
}
func FindUpdatedArticle(link string, isCsBoard bool, lastNotified time.Time) []BoardItems {
	toReturn := []BoardItems{}

	// get all articles in page
	resp, err := soup.Get(link)
	ErrCheck(err)
	docs := soup.HTMLParse(resp)

	var notices []soup.Root
	if isCsBoard {
		notices = docs.Find("div", "class", "list_box").FindAll("li")
	} else {
		notices = docs.Find("div", "class", "bbs-list").FindAll("tr", "class", "notice")
	}
	for _, tr := range notices {
		// get article title and link
		href := tr.Find("div").Find("a").Attrs()["href"]
		title := tr.Find("div").Find("a").Text()

		// get upload date
		uploadDate := ""
		if isCsBoard {
			articleResp, err := soup.Get("https://cs.sogang.ac.kr" + href)
			ErrCheck(err)
			article := soup.HTMLParse(articleResp)
			uploadDate = article.Find("div", "class", "post_info").Find("div", "class", "info").FindAll("span")[1].Text()
		} else {
			articleResp, err := soup.Get("https://sogang.ac.kr" + href)
			ErrCheck(err)
			article := soup.HTMLParse(articleResp)
			uploadDate = article.Find("div", "class", "info").FindAll("div", "class", "unit")[1].Find("span", "class", "value").Text()
		}
		// format upload date
		const format = "2006.01.02 15:04:05"
		parsedUploadDate, err := time.Parse(format, uploadDate)
		ErrCheck(err)

		// check if upload is after lastNotifed
		if parsedUploadDate.After(lastNotified) {
			item := BoardItems{Title: title, Url: href}
			toReturn = append(toReturn, item)
		}
	}
	return toReturn
}
