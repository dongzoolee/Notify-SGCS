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
func SendMessage(receiverInfo updateDB.ChannelWrap, boardName string, title string, url string) {
	var api = slack.New(receiverInfo.TeamToken)
	// var blocks []slack.Block
	// blocks=append(blocks, slack.)
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*"+boardName+"가 업데이트 되었습니다.*", false, false), nil, nil))
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", "학부야", false, false), nil, nil))
	// blocks = append(blocks, slack.NewDividerBlock())
	// blocks = append(blocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "<"+url+"|"+title+">", false, false), nil, nil))
	// api.PostMessage(receiverInfo.ChannelID, slack.MsgOptionBlocks(blocks...))

	var attch []slack.Attachment
	attch = append(attch, slack.Attachment{Pretext: boardName + "가 업데이트 되었습니다.", Title: title, TitleLink: url})
	api.PostMessage(receiverInfo.ChannelID, slack.MsgOptionAttachments(attch...))
}
