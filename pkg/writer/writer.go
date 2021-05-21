package writer

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
)

func WriteDocumentsToFile(file string, collection *mongo.Collection) error {
	var results []bson.M
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	for cursor.Next(context.TODO()) {
		val := bson.M{}
		err := cursor.Decode(&val)
		if err != nil {
			return err
		}
		results = append(results, val)
	}
	bytes, err := json.Marshal(results)
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		logrus.Fatalf("Failed to write file %s: %v", file, err)
	}
	return nil
}
