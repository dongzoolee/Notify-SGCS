package test

import (
	checkupdate "checkUpdate"
	"fmt"
	"testing"
	"updateDB"
)

var boardTypes = []string{
	"cs-main", "cs-underg", "cs-grad", "cs-general", "cs-applause", "cs-job", "scholar", "underg", "general",
}

func TestChannel(t *testing.T) {
	lists := updateDB.FindChannels(boardTypes)
	dummy1 := updateDB.ChannelWrap{
		ChannelID: "C01MSUXUY3E",
		TeamToken: "xoxb-560811458228-1743020221089-lB7Cw0i7ReRCJJ6zzkN8G1kY",
	}
	dummy2 := updateDB.ChannelWrap{
		ChannelID: "C01MSUXUY3E",
		TeamToken: "xoxb-560811458228-1743020221089-lB7Cw0i7ReRCJJ6zzkN8G1kY",
	}
	toCompare := []updateDB.ChannelWrap{dummy1, dummy2}

	if lists[0].ChannelID == toCompare[0].ChannelID {
		if lists[1].ChannelID == toCompare[1].ChannelID {
			return
		}
	}
	t.Fail()
}
func TestArticleCount(t *testing.T) {
	boards := updateDB.FindBoards()
	for _, board := range boards {
		updated := checkupdate.FindUpdatedArticle(board.Link, board.IsCsBoard, board.LastNotified)
		fmt.Print(board.Name, len(updated))
	}
}

/*func TestMessage(t *testing.T) {
	var api = slack.New("xoxb-560811458228-1743020221089-lB7Cw0i7ReRCJJ6zzkN8G1kY")
	var attch []slack.Attachment
	attch = append(attch, slack.Attachment{Pretext: "Notify SGCS Bot 점검 안내", Title: "25일(수) 오후 1시부터 오후 6시까지 점검이 예정되어 있습니다.\r이용에 참고 부탁드립니다."})
	api.PostMessage("U011D7BSBAR", slack.MsgOptionAttachments(attch...))
}*/

/*func TestBroadcast(t *testing.T) {
	check := make(map[string]bool)
	lists := updateDB.GetChannels(boardTypes)
	for _, elem := range lists {
		if check[elem.ChannelID] {
			continue
		}
		check[elem.ChannelID] = true

		var api = slack.New(elem.TeamToken)
		var attch []slack.Attachment
		attch = append(attch, slack.Attachment{Pretext: "Notify SGCS Bot 점검 안내", Title: "25일(수) 오후 1시부터 오후 6시까지 점검이 예정되어 있습니다.\r이용에 참고 부탁드립니다."})
		api.PostMessage(elem.ChannelID, slack.MsgOptionAttachments(attch...))
	}
}*/
