package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoDB *mongo.Database
)

func InitMongoDatabase(uri string, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var opt options.ClientOptions
	opt.SetMaxPoolSize(10)
	opt.SetMinPoolSize(10)

	opt.SetReadPreference(readpref.SecondaryPreferred())
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri), &opt)
	if err != nil {
		fmt.Printf("NEW_MONGO_ERROR %s\n", err.Error())
		return err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("NEW_MONGO_ERROR %s\n", err.Error())
		return err
	}

	MongoDB = mongoClient.Database(dbName)
	return nil
}
