package api

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSettings(cfg *Config) (*Settings, error) {
	lcfg, err := cfg.NewLinBotSetting()
	if err != nil {
		return nil, err
	}

	dbcfg, err := cfg.NewDBSetting()
	if err != nil {
		return nil, err
	}

	svCfg, err := cfg.NewServerSetting()
	if err != nil {
		return nil, err
	}

	return &Settings{
		Line:   lcfg,
		Mongo:  dbcfg,
		Server: svCfg,
	}, nil
}

func (cfg *Config) NewLinBotSetting() (l *LineBotConfig, err error) {
	l = &LineBotConfig{}
	if err := ReadConfig(cfg.LineCfg, l); err != nil {
		return nil, err
	}
	log.Print(l.ChannelAccessToken)

	return l, nil
}

func (cfg *Config) NewDBSetting() (db *MongoDBConfig, err error) {
	db = &MongoDBConfig{}
	if err := ReadConfig(cfg.DBCfg, db); err != nil {
		return nil, err
	}
	log.Print("Finished DB setting")

	return db, nil
}

func (db *MongoDBConfig) ConnectDB() (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%v:%v@127.0.0.1:27017", db.User, db.Pwd)
	log.Print(uri)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB!")

	return client, nil
}

func (cfg *Config) NewServerSetting() (sv *ServerConfig, err error) {
	sv = &ServerConfig{}
	if err := ReadConfig(cfg.ServerCfg, sv); err != nil {
		return nil, err
	}
	log.Println("Server setting in", sv.Port)

	return sv, nil
}

func ReadConfig(cfg string, b interface{}) error {
	vp := viper.New()
	vp.SetConfigFile(cfg)
	vp.AutomaticEnv()

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	if err := vp.Unmarshal(&b); err != nil {
		return err
	}

	return nil
}

func CloseDatabase(client *mongo.Client) {
	if client != nil {
		client.Disconnect(context.Background())
	}
}
