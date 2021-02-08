package main

import (
	"checkUpdate"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"updateDB"
	// "updateDB"
)

type slackBody []struct {
}

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}
func getEnv(key string) string {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	return os.Getenv(key)
}

var api = slack.New(getEnv("BOT_TOKEN"))

func main() {
	h := mux.NewRouter()
	go checkupdate.Init()

	// INIT SLACK API
	signingSecret := getEnv("SLACK_SIGNING_SECRET")
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
	// GET DATA

	// HTTP INIT
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to Main Page")
	})
	h.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		result, err := ioutil.ReadAll(r.Body)
		ErrCheck(err)
		body := strings.Split(string(result), "&")
		channelID, err := url.QueryUnescape(strings.Split(body[3], "=")[1])
		ErrCheck(err)
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
		if cmd[1:] == "on" {
			if updateDB.AddChannel(channelID, text) {
				api.PostMessage(channelID, slack.MsgOptionText(boardType+" 업데이트에 대한 알림을 받습니다.", false))
			} else {
				api.PostMessage(channelID, slack.MsgOptionText("이미 "+boardType+" 업데이트에 대한 알림을 받고 있습니다.", false))
			}
		} else if cmd[1:] == "off" {
			if updateDB.RemoveChannel(channelID, text) {
				api.PostMessage(channelID, slack.MsgOptionText(boardType+" 업데이트에 대한 알림을 더 이상 받지 않습니다.", false))
			} else {
				api.PostMessage(channelID, slack.MsgOptionText("이미 "+boardType+" 업데이트에 대한 알림을 받지 않고 있습니다.", false))
			}
		}
	})
	http.Handle("/", h)

	http.ListenAndServe(":4567", nil)

}
