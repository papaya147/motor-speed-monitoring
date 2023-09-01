package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() (*mongo.Client, error) {
	uriString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		MONGO_USERNAME,
		MONGO_PASSWORD,
		MONGO_URL,
		MONGO_DBNAME,
	)
	opts := options.Client().ApplyURI(uriString)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(MONGO_DBNAME).Collection(MONGO_COLLECTION)
}

type MotorSpeedRecord struct {
	Time       time.Time `bson:"time" json:"time"`
	MotorSpeed float64   `bson:"motor_speed" json:"motor_speed"`
}

func (app *Config) InsertMotorSpeedRecord(record MotorSpeedRecord) error {
	_, err := app.MongoCollection.InsertOne(context.TODO(), record)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) FetchMotorSpeedRecords(start time.Time, end time.Time) ([]MotorSpeedRecord, error) {
	filter := bson.M{
		"time": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := app.MongoCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var records []MotorSpeedRecord
	var record MotorSpeedRecord

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
