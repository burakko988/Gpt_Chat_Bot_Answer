package main

import (
	"chatbot/chat"
	"chatbot/common"
	"chatbot/config"
	"chatbot/database"
	"chatbot/event"
	"chatbot/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// API_ID web socket api id on aws.
const API_ID = "API_ID"

// WS_KEY web socket authentication key.
const WS_KEY = "WS_KEY"

// OPEN_ID_API_KEY open id api key.
const OPEN_ID_API_KEY = "OPEN_ID_API_KEY"

// WS_URL web socket url.
const WS_URL = "wss://" + API_ID + ".execute-api.eu-central-1.amazonaws.com/prod?HeaderAuth=" + WS_KEY + ":400"

const WEB_SOCKET_GOING_AWAY = "websocket: close 1001 (going away): Going away"

var MessageQue = []Que{}

var WebSocket *websocket.Conn

type Que struct {
	Data     model.WebSocketBody
	SendTime time.Time
}

type Cumulative struct {
	Info     string `json:"info"`
	Question string `json:"question"`
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func main() {

	c := cron.New()
	defer c.Stop()

	// Add cron job to run myFunction every 5 minutes
	c.AddFunc("*/1 * * * *", cron1)

	// Start cron
	c.Start()

	// Keep the main function running
	for {
		time.Sleep(1 * time.Second)
	}

}

func cron1() {

	cron := cron.New()

	var logfile *os.File

	if _, err := os.Stat(config.LogFileName); errors.Is(err, os.ErrNotExist) {
		logfile, err = os.Create(config.LogFileName)
		if err != nil {
			log.Fatal(err)
		}
	} else {

		logfile, err = os.OpenFile(config.LogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initlize the log instance.
	common.InitLog(logfile)

	db := database.Connect()

	profileCollection := db.Collection("profiles")
	conversationCollection := db.Collection("conversations")
	messagesCollection := db.Collection("messages")

	bots := database.GetBotAccount(profileCollection)

	botIds := []primitive.ObjectID{}

	for _, b := range bots {
		botIds = append(botIds, b.ID)
	}
	cbids := ClearSocket(botIds)

	arr := []primitive.ObjectID{}

	arr = cbids

	// mn := database.ClearSocket(profileCollection, arr)

	cron.AddFunc("*/1 * * * *", func() {
		database.ClearSocket(profileCollection, arr)
	})

	// need cron every 15 min
	openaiClient := chat.NewService(OPEN_ID_API_KEY)
	start := time.Now()
	common.Sugar.Infow("Trying to connect the web socket", "api_id", API_ID)

	WebSocket, _, err := websocket.DefaultDialer.Dial(WS_URL, nil)

	cron.AddFunc("*/3 * * * *", func() {
		common.Sugar.Infow("CRON job has the time try to send the messages", "api_id", API_ID)
		for _, msg := range MessageQue {

			var j []byte
			if msg.SendTime.Before(time.Now()) {

				start := time.Now()

				elapsed := time.Since(start)

				j, _ = json.Marshal(msg.Data)

				if err != nil {

					common.Sugar.Fatalw("Error on marshaling the event to send message", "api_id", API_ID, "event", "newMessage", "error", err, "duration", elapsed.String())
				}

				if WebSocket != nil {
					err = WebSocket.WriteMessage(1, j)

					if err != nil {

						common.Sugar.Fatalw("Error on sending message", "api_id", API_ID, "event", "newMessage", "error", err, "duration", elapsed.String())
					}

				}

				common.Sugar.Infow("Message sucessfully send", "event", "sendMessage", "conversation", msg.Data.Data.ConversationId, "duration", elapsed.String())
			}
			MessageQue = []Que{}

		}

	})

	cron.Start()

	if err != nil {
		common.Sugar.Fatalw("Error on connecting the socket", "api_id", API_ID, "error", err, "duration", time.Since(start).String())
	}

	common.Sugar.Infow("Successfully connected to the socket", "api_id", API_ID, "duration", time.Since(start).String())

	common.Sugar.Info("Waiting for the incoming messages")

	done := make(chan []byte)

	go func() {
		defer close(done)
		for {
			_, message, err := WebSocket.ReadMessage()
			if err != nil {
				if err.Error() == WEB_SOCKET_GOING_AWAY {
					common.Sugar.Fatalw("Disconnecting from the web socket for inactivity", "api_id", API_ID, "time", time.Now().String())
				}
				common.Sugar.Fatalw("Error on reading the incoming message", "api_id", API_ID, "error", err)
			}

			done <- message
		}
	}()

	for wsEvent := range done {
		eventMap := map[string]interface{}{}

		if err := json.Unmarshal([]byte(wsEvent), &eventMap); err != nil {
			common.Sugar.Fatalw("Error on unmarshaling the event to eventMap struct", "api_id", API_ID, "error", err)
		}

		eventName := eventMap["event"]

		if eventName == "newMessage" {
			common.Sugar.Infow("New message has been received", "event", eventName)

			newMessage := event.MessageSocketEvent{}
			b, err := json.Marshal(eventMap)

			if err != nil {
				common.Sugar.Fatalw("Error on marshaling the event", "api_id", API_ID, "event", eventName, "error", err)
			}

			if err := json.Unmarshal(b, &newMessage); err != nil {
				common.Sugar.Fatalw("Error on unmarshaling the event to new message struct", "api_id", API_ID, "event", eventName, "error", err)
			}

			if !common.IsBotEcho(newMessage.Message.Sender, botIds) {
				common.Sugar.Infow("New message has been received", "content", newMessage.Message.Content)

				start := time.Now()

				conversation := database.GetConversation(conversationCollection, newMessage.ConversationId)
				if conversation != nil {

					botId := conversation.GetBot(botIds)
					// Simple cumulative algorithm get on database last 5 messages cause OPEN-AI earn a token and a little expensive :)
					cumulative := database.GetMessages(messagesCollection, conversation.ID)

					botInfo, _ := database.GetProfileById(profileCollection, *botId)

					if len(botInfo.Sockets) > 0 {

						var items []string

						botName := botInfo.Username

						botGender := botInfo.Meta.Gender
						// HERE AUTOMATICLY ADDED EACH CUMULATIVE CONVERSATION AND BOT ALWAYS GIVE THE TRUE ANSWER.
						info := "AI name:" + botName + "AI gender:" + botGender + "ADD WHAT U WANT :)" + "\n\n"

						strings.Join(items, info)

						for _, c := range *cumulative {
							message := c.Content
							if c.Sender != *botId {
								message = "Human: " + message
							}
							items = append(items, message)
						}

						ri := reverse(items)

						merged := strings.Join(ri, "\n\n")

						if botId != nil {
							common.Sugar.Infow("Bot has been found on the conversation users", "conversation", "botId", botId.Hex(), newMessage.ConversationId.Hex())
							content := openaiClient.GoGpt("Human: "+merged, info)

							sendEvent := model.WebSocketBody{
								Action: "sendmessage",
								Data: model.SendMessageToConversation{
									ConversationId: newMessage.ConversationId.Hex(),
									Content:        strings.ReplaceAll(content.Content, "\n", ""),
									ContentType:    newMessage.Message.Type,
									Caption:        newMessage.Message.Caption,
									Reply:          "",
									IsBot:          botId.Hex(),
									SyncId:         "123",
								},
							}

							prf := AddRandomMinutes(time.Now())

							MessageQue = append(MessageQue, Que{
								Data:     sendEvent,
								SendTime: prf,
							})
							fmt.Println(MessageQue)

						} else {
							common.Sugar.Warnw("Bot did not online", "event", eventName, "conversation", newMessage.ConversationId.Hex(), "duration", time.Since(start))

						}
					} else {

						common.Sugar.Warnw("Bot did not found", "event", eventName, "conversation", newMessage.ConversationId.Hex(), "duration", time.Since(start))
					}
				} else {

					common.Sugar.Warnw("Conversation did not found", "event", eventName, "conversation", newMessage.ConversationId.Hex(), "duration", time.Since(start))
				}

			} else {
				common.Sugar.Warnw("It is new message echo do not take any action", "event", eventName, "conversation", newMessage.ConversationId.Hex())
			}
		} else {
			common.Sugar.Infow("New event has been received", "event", eventName)
		}

	}
}

func AddRandomMinutes(t time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	randomMinutes := rand.Intn(15) + 1
	minutesToAdd := time.Duration(randomMinutes) * time.Minute
	return t.Add(minutesToAdd)
}

func ClearSocket(botIds []primitive.ObjectID) []primitive.ObjectID {

	subset := []primitive.ObjectID{}

	for len(subset) < len(botIds)/2 {

		index := rand.Intn(len(botIds))

		if !contains(subset, botIds[index]) {
			subset = append(subset, botIds[index])
		}
	}
	fmt.Println(subset)
	return subset

}

func contains(arr []primitive.ObjectID, elem primitive.ObjectID) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}
