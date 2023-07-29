package util

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// insertKeywordsToMongoDB inserts the non-repeat keywords
func insertKeywordsToMongoDB(collection *mongo.Collection, keywords []string) error {
	for _, keyword := range keywords {
		filter := bson.D{{"keyword", keyword}}
		update := bson.D{{"$set", bson.D{{"keyword", keyword}}}}
		opt := options.Update().SetUpsert(true)

		_, err := collection.UpdateOne(context.Background(), filter, update, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

// readKeywords reads the file `filename` and gets all those keywords, return a slice of those keywords.
func readKeyword(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var keywords []string
	for scanner.Scan() {
		keyword := strings.TrimSpace(scanner.Text())
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}

func SetKeywords(collection *mongo.Collection, fileName string) error {
	keywords, err := readKeyword(fileName)
	if err != nil {
		return err
	}

	if err := insertKeywordsToMongoDB(collection, keywords); err != nil {
		return err
	}

	return nil
}

// IsSensitive checks if the msg sent by user contained sensitive keywords.
func IsSensitive(text string, collection *mongo.Collection) bool {
	text = strings.ToLower(text)

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Failed to query keywords from MongoDB:", err)
		return false
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var keyword struct{ Keyword string }
		if err := cursor.Decode(&keyword); err != nil {
			log.Println("Failed to decode keyword from MongoDB:", err)
			continue
		}

		if strings.Contains(text, keyword.Keyword) {
			return true
		}
	}

	return false
}
