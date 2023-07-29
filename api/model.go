package api

import (
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	Port          string
	LineBotClient *linebot.Client
	MongoClient   *mongo.Client
	Collection    *mongo.Collection
	KeyCollection *mongo.Collection
}

type Settings struct {
	Line   *LineBotConfig
	Mongo  *MongoDBConfig
	Server *ServerConfig
}

// Config structure
type Config struct {
	LineCfg   string
	ServerCfg string
	DBCfg     string
}

type LineBotConfig struct {
	ChannelSecret      string `mapstructure:"channelSecret"`
	ChannelAccessToken string `mapstructure:"channelAccessToken"`
}

type MongoDBConfig struct {
	User           string `mapstructure:"dbUser"`
	Pwd            string `mapstructure:"dbPwd"`
	DBName         string `mapstructure:"dbName"`
	CollectionName string `mapstructure:"collectionName"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// ReplyMessage structure
type ReplyMessage struct {
	Message string `mapstructure:"message"`
}

type LineEvent struct {
	Id                string    `json:"id" bson:"_id"`
	ReplyToken        string    `json:"replytoken" bson:"replytoken"`
	Type              string    `json:"type" bson:"type"`
	Mode              string    `json:"mode" bson:"mode"`
	Timestamp         time.Time `json:"timestamp" bson:"timestamp"`
	Source            Source    `json:"source" bson:"source"`
	Message           Message   `json:"message" bson:"message"`
	Joined            string    `json:"joined" bson:"joined"`
	Left              string    `json:"left" bson:"left"`
	AccountLink       string    `json:"accountlink" bson:"accountlink"`
	Things            string    `json:"things" bson:"things"`
	Members           string    `json:"members" bson:"members"`
	Unsend            string    `json:"unsend" bson:"unsend"`
	Vedioplaycomplete string    `json:"vedioplaycomplete" bson:"vedioplaycomplete"`
}

type Source struct {
	Type    string `json:"type" bson:"type"`
	UserId  string `json:"userid" bson:"userid"`
	GroupId string `json:"groupid" bson:"groupid"`
	RoomId  string `json:"roomid" bson:"roomid"`
}

type Emojis struct {
	Index     int    `json:"index" bson:"index"`
	Length    int    `json:"length" bson:"length"`
	ProductId string `json:"productid" bson:"productid"`
	EmojiId   string `json:"emojiid" bson:"emojiid"`
}

type Message struct {
	Id     string   `json:"id" bson:"id"`
	Text   string   `json:"text" bson:"text"`
	Emojis []Emojis `json:"emojis" bson:"emojis"`
	Metion string   `json:"metion" bson:"metion"`
}
