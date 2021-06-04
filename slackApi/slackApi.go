package SlackApi

import (
	"log"
	"os"
	"updateDB"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
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
func SendMsg(receiverInfo updateDB.ChannelWrap, boardType string, title string, url string) {
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
	var api = slack.New(receiverInfo.TeamToken)
	// var api = slack.New("xoxb-560811458228-1743020221089-bdbPRFN8B0VH5fTpoo5BSR3J")
	// var blocks []slack.Block
	// blocks=append(blocks, slack.)
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*"+boardType+"가 업데이트 되었습니다.*", false, false), nil, nil))
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", "학부야", false, false), nil, nil))
	// blocks = append(blocks, slack.NewDividerBlock())
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "<"+url+"|"+title+">", false, false), nil, nil))
	// api.PostMessage(receiverInfo.ChannelID, slack.MsgOptionBlocks(blocks...))
	// api.PostMessage("U011D7BSBAR", slack.MsgOptionBlocks(blocks...))
	// api.PostMessage("U011D7BSBAR", slack.MsgOptionText("학부야", true))

	var attch []slack.Attachment
	attch = append(attch, slack.Attachment{Pretext: boardType + "가 업데이트 되었습니다.", Title: title, TitleLink: url})
	// api.PostMessage("U011D7BSBAR", slack.MsgOptionAttachments(attch...))
	api.PostMessage(receiverInfo.ChannelID, slack.MsgOptionAttachments(attch...))
}
