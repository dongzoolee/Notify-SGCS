package SlackApi

import (
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"log"
	"os"
)

func GetEnv(key string) string {
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

var api = slack.New(GetEnv("BOT_TOKEN"))

func SendMsg(channel string, boardType string, title string, url string) {
	if boardType == "main" {
		boardType = "주요공지"
	} else if boardType == "underg" {
		boardType = "학부공지"
	} else if boardType == "grad" {
		boardType = "대학원공지"
	} else if boardType == "general" {
		boardType = "일반공지"
	} else if boardType == "job" {
		boardType = "취업/인턴십공지"
	} else if boardType == "sgcs" {
		boardType = "학과소식"
	}
	var blocks []slack.Block
	blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*"+boardType+"가 업데이트 되었습니다.*\n> <"+url+"|"+title+">", false, false), nil, nil))
	api.PostMessage(channel, slack.MsgOptionBlocks(blocks...))
}
