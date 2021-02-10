package main

import (
	checkupdate "checkUpdate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"updateDB"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	// "updateDB"
)

type slackBody []struct {
}

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}
func Getenv(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	return os.Getenv(key)
}

var api = slack.New(Getenv("BOT_TOKEN"))

func main() {
	h := mux.NewRouter()
	go checkupdate.Init()

	// INIT SLACK API
	signingSecret := Getenv("SLACK_SIGNING_SECRET")
	h.HandleFunc("/event-endpoint", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				fmt.Println(ev.Channel, ev.Text)
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			case *slackevents.MessageEvent:
				if ev.ClientMsgID == "" {
					return
				}
				if ev.Text == "on main" {
					api.PostMessage(ev.Channel, slack.MsgOptionText("주요공지를 알림 받습니다.", false))
				}
				// default:
			}
		}
	})

	// SETTING INSTALL PAGE
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://csnotice.soga.ng/install", 301)
	})
	h.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://slack.com/oauth/v2/authorize?client_id="+Getenv("SLACK_CLIENT_ID")+"&scope=chat%3Awrite+commands+chat%3Awrite.public+im%3Awrite", 301)
	})
	h.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		query := strings.Split(r.URL.RawQuery, "?")
		code := strings.Split(query[0], "&")[0]
		queryOfCode := strings.Split(code, "=")[1]
		resp, err := http.Get("https://slack.com/api/oauth.v2.access?client_id=" + Getenv("SLACK_CLIENT_ID") + "&client_secret=" + Getenv("SLACK_CLIENT_SECRET") + "&code=" + queryOfCode)
		ErrCheck(err)
		defer resp.Body.Close()

		result, err := ioutil.ReadAll(resp.Body)
		ErrCheck(err)
		var res map[string]interface{}
		json.Unmarshal([]byte(result), &res)
		fmt.Println(res)
		updateDB.SetTeamToken(res["team"].(map[string]interface{})["id"].(string), res["access_token"].(string))
		http.Redirect(w, r, "/installsuccess", 301)
	})
	h.HandleFunc("/installsuccess", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Successfully Installed in your workspace"))
	})

	// RECEIVE SLASH EVENTS
	h.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		result, err := ioutil.ReadAll(r.Body)
		ErrCheck(err)
		body := strings.Split(string(result), "&")
		// fmt.Println(body)
		if len(body) == 0 {
			return
		}

		// DM의 경우, 사용자 ID로 메세지를 보내야하고,
		var channelID string
		if strings.Split(body[4], "=")[1] == "directmessage" {
			channelID = strings.Split(body[5], "=")[1]
		}else{ // 일반 채널은 채널 ID로 메세지를 보내야 합니다.
			channelID = strings.Split(body[3], "=")[1]
		}
		cmd, err := url.QueryUnescape(strings.Split(body[7], "=")[1])
		ErrCheck(err)
		text, err := url.QueryUnescape(strings.Split(body[8], "=")[1])
		ErrCheck(err)
		var boardType string
		if text == "main" {
			boardType = "주요공지"
		} else if text == "underg" {
			boardType = "학부공지"
		} else if text == "grad" {
			boardType = "대학원공지"
		} else if text == "general" {
			boardType = "일반공지"
		} else if text == "job" {
			boardType = "취업/인턴십공지"
		} else if text == "sgcs" {
			boardType = "학과소식"
		}
		// fmt.Println(updateDB.GetTeamToken(strings.Split(body[1], "=")[1]))
		var api = slack.New(updateDB.GetTeamToken(strings.Split(body[1], "=")[1]))
		if cmd[1:] == "on" {
			// api.PostMessage(channelID, slack.MsgOptionText("유지보수 중입니다 잠시 후 시도해주세요", false))
			// return
			if updateDB.AddChannel(updateDB.GetTeamToken(strings.Split(body[1], "=")[1]), strings.Split(body[1], "=")[1], channelID, text) {
				// resp, err := http.PostForm("https://slack.com/api/chat.postMessage", url.Values{"token": {"xoxb-1691653951013-1729473312866-18XAO49tpUhwGWVgvy1IrLsa"}, "channel": {channelID}, "blocks": {`[{"type": "section", "text": {"type": "plain_text", "text": "Hello world"}}]`}})
				// ErrCheck(err)
				// respBody, err := ioutil.ReadAll(resp.Body)
				// ErrCheck(err)
				// if err == nil {
				// 	str := string(respBody)
				// 	fmt.Println(str)
				// }
				// fmt.Println(channelID)
				_, _, err = api.PostMessage(channelID, slack.MsgOptionText(boardType+" 업데이트에 대한 알림을 받습니다.", false))
				ErrCheck(err)
			} else {
				// fmt.Println(channelID)
				_, _, err := api.PostMessage(channelID, slack.MsgOptionText("이미 "+boardType+" 업데이트에 대한 알림을 받고 있습니다.", false))
				ErrCheck(err)
			}
		} else if cmd[1:] == "off" {
			// api.PostMessage(channelID, slack.MsgOptionText("유지보수 중입니다 잠시 후 시도해주세요", false))
			// return
			if updateDB.RemoveChannel(channelID, text) {
				api.PostMessage(channelID, slack.MsgOptionText(boardType+" 업데이트에 대한 알림을 더 이상 받지 않습니다.", false))
				ErrCheck(err)
			} else {
				api.PostMessage(channelID, slack.MsgOptionText("이미 "+boardType+" 업데이트에 대한 알림을 받지 않고 있습니다.", false))
				ErrCheck(err)
			}
		}
	})
	http.Handle("/", h)

	http.ListenAndServe(":4567", nil)

}
