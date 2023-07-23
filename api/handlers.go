package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *API) receiveHandler(ctx *gin.Context) {
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

func (a *API) getAllMessages(ctx *gin.Context) {
	var linemessage LineEvent
	var results []LineEvent

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set a cursor
	cursor, err := a.Collection.Find(c, bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// loop through all results and save in results slice
	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&linemessage); err != nil {
			log.Fatal(err)
		}
		results = append(results, linemessage)
	}

	ctx.JSON(http.StatusOK, results)
}
