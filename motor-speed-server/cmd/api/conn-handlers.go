package main

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

var count = 0

func ConnectToMqtt() mqtt.Client {
	for {
		client, err := NewMqttClient()
		if err != nil {
			log.Println("mqtt not yet ready...")
			count++
		} else {
			log.Println("mqtt ready")
			return client
		}

		if count > 10 {
			log.Println(err)
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func ConnectToMongo() *mongo.Client {
	for {
		client, err := NewMongoClient()
		if err != nil {
			log.Println("mongo not yet ready...")
			count++
		} else {
			log.Println("mongo ready")
			return client
		}

		if count > 10 {
			log.Println(err)
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
