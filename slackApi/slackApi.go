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
func SendMessage(receiverInfo updateDB.ChannelWrap, boardName string, title string, url string) slack.Attachment {
	var api = slack.New(receiverInfo.TeamToken)

	var attachment = slack.Attachment{Pretext: boardName + "가 업데이트 되었습니다.", Title: title, TitleLink: url}
	_, _, err := api.PostMessage(receiverInfo.ChannelID, slack.MsgOptionAttachments(
		[]slack.Attachment{attachment}...,
	))

	if err != nil {
		log.Panic(err)
	}

	return attachment
}
