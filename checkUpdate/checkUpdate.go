package checkupdate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	SlackApi "slackApi"
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
func Init() {
	t := time.NewTicker(30 * time.Minute)
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
	var channelIDList = updateDB.GetChannels(boardType)
	chk, ret := CmpPKID(boardType)
	if chk {
		for _, receiverInfo := range channelIDList {
			for _, val := range ret {
				SlackApi.SendMsg(receiverInfo, boardType, val.Title, val.Url)
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
	var oldStrongTopPostID string
	var oldGeneralTopPostID string
	// Set Board Type's ID
	var boardID string
	if boardType == "main" {
		boardID = "1905"
		oldStrongTopPostID = unMshedD1.Strong.Main
		oldGeneralTopPostID = unMshedD1.General.Main
	} else if boardType == "underg" {
		boardID = "1745"
		oldStrongTopPostID = unMshedD1.Strong.Underg
		oldGeneralTopPostID = unMshedD1.General.Underg
	} else if boardType == "grad" {
		boardID = "1747"
		oldStrongTopPostID = unMshedD1.Strong.Grad
		oldGeneralTopPostID = unMshedD1.General.Grad
	} else if boardType == "general" {
		boardID = "1746"
		oldStrongTopPostID = unMshedD1.Strong.General
		oldGeneralTopPostID = unMshedD1.General.General
	} else if boardType == "job" {
		boardID = "1748"
		oldStrongTopPostID = unMshedD1.Strong.Job
		oldGeneralTopPostID = unMshedD1.General.Job
	} else if boardType == "sgcs" {
		boardID = "1749"
		oldStrongTopPostID = unMshedD1.Strong.Sgcs
		oldGeneralTopPostID = unMshedD1.General.Sgcs
	}
	// fmt.Print(time.Now())
	resp, err := soup.Get("https://cs.sogang.ac.kr/front/cmsboardlist.do?siteId=cs&bbsConfigFK=" + boardID)
	ErrCheck(err)
	// fmt.Println("SUCCESS")
	fmt.Println("boardID: " + boardID + ", oldStrID: " + oldStrongTopPostID + ", oldGenID: " + oldGeneralTopPostID)
	doc := soup.HTMLParse(resp)
	lis := doc.Find("div", "class", "list_box").FindAll("li")

	var newStrongTopPostID = oldStrongTopPostID
	var newGeneralTopPostID = oldGeneralTopPostID
	var ret []BoardItems
	var isUpdated bool = false
	var case1 bool = false
	var case2 bool = false
	var isCase2TopIdx bool = true
	var generalPostCnt int = 0
	for idx, li := range lis {
		href := li.Find("div").Find("a").Attrs()["href"]
		title := li.Find("div").Find("a").Text()
		strong := li.Find("div").Find("a").Find("strong")
		splitedHref := strings.Split(href, "&")
		postID := strings.Split(splitedHref[len(splitedHref)-1], "=")[1]
		if idx == 0 {
			// [공지]라면
			if strong.Error == nil {
				// 업데이트가 됐다
				if oldStrongTopPostID != postID {
					tmp := new(BoardItems)
					tmp.Title = title
					tmp.Url = "https://cs.sogang.ac.kr" + href
					ret = append(ret, *tmp)

					fmt.Println("StrongPost Updated. Appending Queue : " + postID)

					isUpdated = true
					newStrongTopPostID = postID
					case1 = true
				} else { // 업데이트 안됐으니 일반공지 돌자
					fmt.Println("StrongPost Nothing Updated. Starting Case2 : " + postID)

					case2 = true
				}
			} else { // 일반공지라면
				//업데이트 됐다
				if oldGeneralTopPostID != postID {
					tmp := new(BoardItems)
					tmp.Title = title
					tmp.Url = "https://cs.sogang.ac.kr" + href
					ret = append(ret, *tmp)

					fmt.Println("GeneralPost Updated. Appending Queue : " + postID)

					isUpdated = true
					isCase2TopIdx = false
					newGeneralTopPostID = postID
					case2 = true
				} else { // 업데이트 안됐다
					// oldStrongTopPostID = "-1" 굳이?
					fmt.Println("GeneralPost Nothing Updated. Stopping Loop : " + postID)

					break
				}
			}
			continue
		}
		// fmt.Println(strong.Error)
		// [공지]라면
		if strong.Error == nil {
			// [공지]가 업데이트 됐었음
			if case1 {
				// 이번 것도 업데이트 된 건지 확인
				if oldStrongTopPostID != postID {
					fmt.Println("StrongPost Still Being Updated. Appending Queue : " + postID)
					tmp := new(BoardItems)
					tmp.Title = title
					tmp.Url = "https://cs.sogang.ac.kr" + href
					ret = append(ret, *tmp)
				} else { // 안됐다면 case1 종료
					fmt.Println("StrongPost Nothing More To Append. Closing Case1 And Start Case2 : " + postID)

					case1 = false
					case2 = true
				}
			} else if case2 { // 일반 공지를 돌아야 함
			}
		} else { // 일반 공지
			generalPostCnt++ // 일반 공지 1페이지에 몇개 있는지

			// 일단 일반 공지 첫번째 idx를 확인
			if isCase2TopIdx {
				// 첫번째 공지가 업데이트된 상태일 경우
				if oldGeneralTopPostID != postID {
					tmp := new(BoardItems)
					tmp.Title = title
					tmp.Url = "https://cs.sogang.ac.kr" + href
					ret = append(ret, *tmp)

					fmt.Println("GeneralPost Still Being Updated. Appending Queue : " + postID)

					isUpdated = true
					newGeneralTopPostID = postID
				} else { // 첫번째 공지가 업데이트 x일 경우
					fmt.Println("GeneralPost Nothing To Append. Stopping Loop : " + postID)
					break
				}
				isCase2TopIdx = false
			} else {
				// 아직도 업데이트된 게시글일 경우
				if oldGeneralTopPostID != postID {
					tmp := new(BoardItems)
					tmp.Title = title
					tmp.Url = "https://cs.sogang.ac.kr" + href
					ret = append(ret, *tmp)

					fmt.Println("GeneralPost Still Being Updated. Appending Queue : " + postID)

				} else { // 드디어 기존 공지를 만났을 경우
					fmt.Println("GeneralPost Nothing More To Append. Stopping Loop : " + postID)

					break
				}
			}
		}
	}

	if case1 && generalPostCnt == 0 { // 2페이지로 일반공지가 넘어갔을 경우
		fmt.Println("페이지 넘어감 ㅜ")
		return false, ret
	}
	if !isUpdated {
		return false, ret
	}

	// reflect to json and save file
	if boardType == "main" {
		unMshedD1.Strong.Main = newStrongTopPostID
		unMshedD1.General.Main = newGeneralTopPostID
	} else if boardType == "underg" {
		unMshedD1.Strong.Underg = newStrongTopPostID
		unMshedD1.General.Underg = newGeneralTopPostID
	} else if boardType == "grad" {
		unMshedD1.Strong.Grad = newStrongTopPostID
		unMshedD1.General.Grad = newGeneralTopPostID
	} else if boardType == "general" {
		unMshedD1.Strong.General = newStrongTopPostID
		unMshedD1.General.General = newGeneralTopPostID
	} else if boardType == "job" {
		unMshedD1.Strong.Job = newStrongTopPostID
		unMshedD1.General.Job = newGeneralTopPostID
	} else if boardType == "sgcs" {
		unMshedD1.Strong.Sgcs = newStrongTopPostID
		unMshedD1.General.Sgcs = newGeneralTopPostID
	}

	mshedD1, err := json.Marshal(unMshedD1)
	ErrCheck(err)
	err = ioutil.WriteFile("./checkUpdate/pkid.json", mshedD1, 0644)
	ErrCheck(err)
	return true, ret
}
