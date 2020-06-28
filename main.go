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
	"bytes"
	"strings"
)

var (
	slackClient *slack.Client
)

type SearchResults struct {
	Items []Item `json:items` 
}

type Item struct {
	Link string `json:link`
	Snippet string `json:snippet`
	Title string `json:title`
}

type TextInfo struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Block struct {
	Type string `json:"type"`
	Text TextInfo `json:"text"`
	BlockId string `json:"block_id"`
}

type Payload struct {
	Blocks []Block `json:"blocks"`
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
	searchAnswer(ev)
}

func replyToUser(jsonMessage []byte) {	

	resp, err := http.Post(os.Getenv("WEB_HOOK"), "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, getErr := ioutil.ReadAll(resp.Body)
	if getErr != nil {
		log.Fatalln(getErr)
	}
	log.Println(body)
}

func searchAnswer(ev *slack.MessageEvent) {
		url := "https://www.googleapis.com/customsearch/v1?key=AIzaSyD8QNzBdjzt3ZNEbGTz4P1rSAnvDPtbrUU&cx=005033773481765961543:gti8czyzyrw&num=3&q=golang"
		googleClient := http.Client{
			Timeout: time.Second * 3, // Maximum of 3 secs
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("User-Agent", "Isacc Hernandez")
		res, getErr := googleClient.Do(req)
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
		apiMessage(body)
}

func apiMessage(jsonRaw []byte) {
	structure := SearchResults{}
	jsonErr := json.Unmarshal(jsonRaw, &structure)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
  dataBinding(structure)
}

func dataBinding(data SearchResults) {
	payload := new(Payload)
	Item1 := data.Items[0]
	Item2 := data.Items[1]
	Item3 := data.Items[2]
	textBlock1 := fmt.Sprintf("*<%s|%s>*\n>_%s_", Item1.Link, Item1.Title, strings.Replace(Item1.Snippet, "\n", " ", -1))
	textBlock2 := fmt.Sprintf("*<%s|%s>*\n>_%s_", Item2.Link, Item2.Title, strings.Replace(Item2.Snippet, "\n", " ", -1))
	textBlock3 := fmt.Sprintf("*<%s|%s>*\n>_%s_", Item3.Link, Item3.Title, strings.Replace(Item3.Snippet, "\n", " ", -1))
	payload.Blocks = []Block{
		Block {
			Type: "section",   
			Text: TextInfo{"mrkdwn", textBlock1}, //Si existe error posiblemente sea porque textBlock no sea un string
			BlockId: "text0",
		},
		Block {
			Type: "section",   
			Text: TextInfo{"mrkdwn", textBlock2}, //Si existe error posiblemente sea porque textBlock no sea un string
			BlockId: "text1",
		},
		Block {
			Type: "section",   
			Text: TextInfo{"mrkdwn", textBlock3}, //Si existe error posiblemente sea porque textBlock no sea un string
			BlockId: "text2",
		},
	}
	jsonMessage, recieveErr := json.Marshal(payload)
	if recieveErr != nil {
		log.Fatalln(recieveErr)
	}
	replyToUser(jsonMessage)
}