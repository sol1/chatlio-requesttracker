package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LunaNode/rtgo"
)

//store configuration data containing user details
type rtConfig struct {
	URL      string
	Username string
	Password string
	Queue    string
}

var configFileName string
var config rtConfig
var rtConn *rtgo.RT

//initialize application flags, load user details from config file.
func init() {
	flag.StringVar(&configFileName, "config", "config.json", "Path to config file.")
	flag.Parse()

	f, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
	}
}

//Main runs http server if appropriate flags specified
func main() {
	if config.URL == "" {
		fmt.Println("Looks like your configuration file is incomplete.")
		fmt.Println("An RT URL MUST be supplied.")
	} else {
		fmt.Println("Server running.")
		http.HandleFunc("/transcript", transcriptHandler)
		http.HandleFunc("/offline", offlineHandler)
		http.ListenAndServe(":8080", nil)
	}
}

//Transcript data structure for decoded JSON, includes requried JSON tags
type Transcript struct {
	ConversationUUID string    `json:"conversationUUID"`
	ChannelID        string    `json:"channelId"`
	SlackTeamID      string    `json:"slackTeamID"`
	ChannelName      string    `json:"channelName"`
	StartedAt        string    `json:"startedAt"`
	EndedAt          string    `json:"endedAt"`
	AgentEmail       string    `json:"agentEmail"`
	AgentName        string    `json:"agentName"`
	VisitorEmail     string    `json:"visitorEmail"`
	ReplyToEmail     string    `json:"replyToEmail"`
	VisitorName      string    `json:"visitorName"`
	VisitorTraits    string    `json:"visitorTraits"`
	ArchiveURL       string    `json:"archiveUrl"`
	Messages         []Message `json:"messages"`
	TextBody         string    `json:"textBody"`
	HTMLBody         string    `json:"htmlBody"`
}

//Message data structure holds all transcript messages
type Message struct {
	Text     string `json:"text"`
	TS       string `json:"ts"`
	SentAt   string `json:"sentAt"`
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

//Handles incoming requests under the /transcript path
func transcriptHandler(w http.ResponseWriter, r *http.Request) {
	//Convert post data from JSON to a readable go format
	var transcript Transcript
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&transcript)
	if err != nil {
		fmt.Printf("Error decoding json body: %s\n", err)
		http.Error(w, "Error decoding json body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// if message is present create RT ticket of the transcript.
	if len(transcript.Messages) > 0 {
		rtConn = rtgo.NewRT(config.URL, config.Username, config.Password)
		_, err = rtConn.CreateTicket(
			config.Queue,
			transcript.VisitorEmail,
			"Live support transcript - "+transcript.VisitorName,
			transcript.TextBody)

		if err != nil {
			fmt.Printf("Error creating ticket in RT: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

//OfflineMessage data structure
type OfflineMessage struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	EmailPretty string `json:"emailPretty"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

//Handle offline chatlio messages
func offlineHandler(w http.ResponseWriter, r *http.Request) {
	//Convert post data from JSON to a readable go format
	var message OfflineMessage
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&message)
	if err != nil {
		fmt.Printf("Error decoding json body: %s\n", err)
		http.Error(w, "Error decoding json body: "+err.Error(), http.StatusBadRequest)
		return
	}

	//Create RT ticket of message
	if message.Message != "" {
		rtConn = rtgo.NewRT(config.URL, config.Username, config.Password)
		_, err = rtConn.CreateTicket(config.Queue, message.Email, message.Subject, message.Message)
		if err != nil {
			fmt.Printf("Error creating ticket in RT: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
