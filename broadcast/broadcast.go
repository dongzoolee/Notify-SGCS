package broadcast

import (
	"updateDB"

	"github.com/slack-go/slack"
)

var boardTypes = []string{
	"general", "grad", "job", "main", "sgcs", "underg",
}

func Broadcast() {
	check := make(map[string]bool)
	lists := updateDB.FindChannels(boardTypes)
	for _, elem := range lists {
		if check[elem.ChannelID] {
			continue
		}
		check[elem.ChannelID] = true

		var api = slack.New(elem.TeamToken)
		var attch []slack.Attachment
		attch = append(attch, slack.Attachment{Pretext: "Notify SGCS Bot 점검 안내", Title: "25일(수) 오후 1시부터 오후 6시까지 점검이 예정되어 있습니다.\r이용에 참고 부탁드립니다."})
		_, _, err := api.PostMessage(elem.ChannelID, slack.MsgOptionAttachments(attch...))
		updateDB.ErrCheck(err)
	}
}
