package api

import (
	"fmt"
	"log"

	"m800_homework/api/util"

	"github.com/line/line-bot-sdk-go/linebot"
)

func NewServer(cfg *Config) (*API, error) {
	settings, err := NewSettings(cfg)
	if err != nil {
		return nil, err
	}

	// Load the 'channelaccesstoken and 'channelsecret' to connect the linebot.
	line, err := linebot.New(settings.Line.ChannelSecret, settings.Line.ChannelAccessToken)
	if err != nil {
		return nil, err
	}

	db, err := settings.Mongo.ConnectDB()
	if err != nil {
		log.Print("Failed to connect to db,", err)
		return nil, err
	}
	collection := db.Database(settings.Mongo.DBName).Collection(settings.Mongo.CollectionName)
	keywordCollection := db.Database(settings.Mongo.DBName).Collection("words")
	if err := util.SetKeywords(keywordCollection, "./config/keywords.txt"); err != nil {
		return nil, err
	}
	log.Print(settings.Mongo.DBName, settings.Mongo.CollectionName)

	return &API{
		Collection:    collection,
		KeyCollection: keywordCollection,
		Port:          settings.Server.Port,
		LineBotClient: line,
		MongoClient:   db,
	}, nil
}

func (a *API) Run() error {
	log.Print("Run server")
	r := NewRouter(a)
	r.Run(fmt.Sprint(":", a.Port))

	return nil
}
