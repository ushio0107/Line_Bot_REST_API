package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
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
