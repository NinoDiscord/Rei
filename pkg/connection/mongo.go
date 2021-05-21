package connection

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func CreateMongoClient(uri string) (client *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("Rei"))
	if err != nil {
		log.Fatal(err)
	}
	
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.Fatalf("connection couldn't be established: %v", err)
		return nil, err
	}

	println("Connected to MongoDB with URI", uri)
	return client, nil
}
