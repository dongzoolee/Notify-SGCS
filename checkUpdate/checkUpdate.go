package checkupdate

import (
	"encoding/json"
	"io/ioutil"
	"slackApi"
	"strings"
	"time"
	"updateDB"

	"github.com/anaskhan96/soup"
)

type BoardItems struct {
	Title string
	Url   string
}
type PKID struct {
	Main    string
	Underg  string
	Grad    string
	General string
	Job     string
	Sgcs    string
}

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}
func Init() {
	t := time.NewTicker(10 * time.Second)
	for range t.C {
		MapDatas("main")
		MapDatas("underg")
		MapDatas("grad")
		MapDatas("general")
		MapDatas("job")
		MapDatas("sgcs")
	}
}
func MapDatas(boardType string) {
	var channelIDList []string = updateDB.GetChannels(boardType)

	chk, ret := CmpPKID(boardType)
	if chk {
		for _, id := range channelIDList {
			for _, val := range ret {
				SlackApi.SendMsg(id, boardType, val.Title, val.Url)
			}
		}
	}
}
func CmpPKID(boardType string) (bool, []BoardItems) {
	d1, err := ioutil.ReadFile("./checkUpdate/pkid.json")
	ErrCheck(err)
	// fmt.Println(string(d1))

	unMshedD1 := new(PKID)
	json.Unmarshal([]byte(d1), &unMshedD1)
	// Get Latest crawled post's ID
	var oldTopPostID string
	// Set Board Type's ID
	var boardID string
	if boardType == "main" {
		boardID = "1905"
		oldTopPostID = unMshedD1.Main
	} else if boardType == "underg" {
		boardID = "1745"
		oldTopPostID = unMshedD1.Underg
	} else if boardType == "grad" {
		boardID = "1747"
		oldTopPostID = unMshedD1.Grad
	} else if boardType == "general" {
		boardID = "1746"
		oldTopPostID = unMshedD1.General
	} else if boardType == "job" {
		boardID = "1748"
		oldTopPostID = unMshedD1.Job
	} else if boardType == "sgcs" {
		boardID = "1749"
		oldTopPostID = unMshedD1.Sgcs
	}

	resp, err := soup.Get("https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=" + boardID)
	ErrCheck(err)
	doc := soup.HTMLParse(resp)
	lis := doc.Find("div", "class", "list_box").FindAll("li")

	var newTopPostID string
	var ret []BoardItems
	var isUpdated bool = false
	for idx, li := range lis {
		href := li.Find("div").Find("a").Attrs()["href"]
		splitedHref := strings.Split(href, "&")
		postID := strings.Split(splitedHref[len(splitedHref)-1], "=")[1]
		if idx == 0 {
			newTopPostID = postID
		}

		if postID != oldTopPostID { // 업데이트된 게시글을 배열에 저장
			tmp := new(BoardItems)
			tmp.Title = li.Find("div").Find("a").Text()
			tmp.Url = "https://cs.sogang.ac.kr" + li.Find("div").Find("a").Attrs()["href"]
			ret = append(ret, *tmp)
			isUpdated = true
		} else if idx == 0 && postID == oldTopPostID { // 제일 위 게시글이 업데이트 되지 않았다면
			break
		} else { // 업데이트된 게시글들을 잘 찾다가 기존의 게시물을 만난다면
			break
		}
	}

	if !isUpdated {
		return false, ret
	}
	// reflect to json and save file
	if boardType == "main" {
		unMshedD1.Main = newTopPostID
	} else if boardType == "underg" {
		unMshedD1.Underg = newTopPostID
	} else if boardType == "grad" {
		unMshedD1.Grad = newTopPostID
	} else if boardType == "general" {
		unMshedD1.General = newTopPostID
	} else if boardType == "job" {
		unMshedD1.Job = newTopPostID
	} else if boardType == "sgcs" {
		unMshedD1.Sgcs = newTopPostID
	}

	mshedD1, err := json.Marshal(unMshedD1)
	ErrCheck(err)
	err = ioutil.WriteFile("./checkUpdate/pkid.json", mshedD1, 0644)
	ErrCheck(err)
	return true, ret
}
