package main

import(
	"os"
	"fmt"
	"log"
	"github.com/slack-go/slack"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
)

var (
	slackClient *slack.Client
)

type searchResults struct {
	Kind string "json:kind"
}

func main() {
	slackClient = slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if len(ev.BotID) == 0 {
				go handleMessage(ev)
			}
		}
	}
}

func handleMessage(ev *slack.MessageEvent) {
	fmt.Printf("%v\n", ev)
	replyToUser(ev)
}

func replyToUser(ev *slack.MessageEvent) {
	slackClient.PostMessage(ev.Channel, slack.MsgOptionText("Hello, User! How can I help you?", false))
	searchAnswer()
}

func searchAnswer() {
		url := "https://www.googleapis.com/customsearch/v1?key=AIzaSyD8QNzBdjzt3ZNEbGTz4P1rSAnvDPtbrUU&cx=005033773481765961543:gti8czyzyrw&num=3&q=golang"
		spaceClient := http.Client{
			Timeout: time.Second * 3, // Maximum of 3 secs
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "Isacc Hernandez")
		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}
		if res.Body != nil {
				defer res.Body.Close()
			}
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		search1 := searchResults{}
		jsonErr := json.Unmarshal(body, &search1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		fmt.Println(search1.Kind)
}