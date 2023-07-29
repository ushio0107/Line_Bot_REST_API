package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"m800_homework/api/util"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

// receiveHandler is the API handler, which reply the user if it's received any message,
// and save the user message to the DB.
func (a *API) receiveHandler(ctx *gin.Context) {
	events, err := a.LineBotClient.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.Writer.WriteHeader(400)
		} else {
			log.Print(err)
			ctx.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if _, err := a.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Get Message")).Do(); err != nil {
				log.Print("Failed to reply message: ", err)
			}

			_, err := a.Collection.InsertOne(ctx, event)
			if err != nil {
				log.Print("Failed to add message to db: ", err)
			}
		}
	}
}

// broadcastMessage broadcasts message to all users.
func (a *API) broadcastMessage(ctx *gin.Context) {
	var message ReplyMessage
	json.NewDecoder(ctx.Request.Body).Decode(&message)
	if err := ctx.Bind(&message); err != nil {
		log.Fatal(err)
	}

	results, err := a.Collection.Distinct(context.TODO(), "source.userid", bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		if val, ok := result.(string); ok {
			if _, err := a.LineBotClient.PushMessage(val, linebot.NewTextMessage(message.Message)).Do(); err != nil {
				fmt.Print("Failed to sent message ", message.Message)
				log.Fatal(err)
			}
		}
	}
}

// getAllMessages gets all messages stored in db.
func (a *API) getAllMessages(ctx *gin.Context) {
	var event LineEvent
	var results []LineEvent

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := a.Collection.Find(c, bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&event); err != nil {
			log.Fatal(err)
		}
		results = append(results, event)
	}

	ctx.JSON(http.StatusOK, results)
}

func (a *API) filterMessage(ctx *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		ctx.Abort()
		return
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	events, err := a.LineBotClient.ParseRequest(ctx.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.Writer.WriteHeader(400)
		} else {
			ctx.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			text := message.Text
			if util.IsSensitive(text, a.KeyCollection) {
				if _, err := a.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Sensitive message")).Do(); err != nil {
					log.Print("Failed to reply message: ", err)
				}
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Sensitive message"})
				ctx.Abort()
				return
			}
		}
	}
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	ctx.Next()
}
