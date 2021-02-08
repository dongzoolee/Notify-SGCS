package slackApi

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"log"
	"encoding/json"
	"github.com/joho/godotenv"
	"os"
)


func getEnv(key string) string {
	err := godotenv.Load("../.env")
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
func init() {
	
}

var api = slack.New(getEnv("BOT_TOKEN"))
func sendMsg() {
	signingSecret := getEnv("SLACK_SIGNING_SECRET")
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
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
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}
	})
}
