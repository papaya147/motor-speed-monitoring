package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

type Data struct {
	MotorSpeed float64 `json:"motor_speed"`
}

const (
	MQTT_URL      = "ecfd1263e29f4248b00d94e4735d8ffb.s2.eu.hivemq.cloud"
	MQTT_PORT     = "8883"
	MQTT_USERNAME = "back-end"
	MQTT_PASSWORD = "3Motorad"
	MQTT_TOPIC    = "motor-speed"

	MONGO_URL        = "motor-speed-test.dd2oemr.mongodb.net"
	MONGO_PORT       = "27017"
	MONGO_USERNAME   = "back-end"
	MONGO_PASSWORD   = "1H3dtMfeHIsDgAHR"
	MONGO_DBNAME     = "motor-speed"
	MONGO_COLLECTION = "motor-speed"

	WEB_PORT = "8080"
)

type Config struct {
	MqttClient      mqtt.Client
	MongoClient     *mongo.Client
	MongoCollection *mongo.Collection
}

func main() {
	mqttClient := ConnectToMqtt()
	mongoClient := ConnectToMongo()
	mongoCollection := NewMongoCollection(mongoClient)

	app := Config{
		MqttClient:      mqttClient,
		MongoClient:     mongoClient,
		MongoCollection: mongoCollection,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	messageChannel := make(chan Data)
	err := app.Subscribe(messageChannel)
	if err != nil {
		log.Panicln("mqtt topic unable to connect, ", err)
	}
	defer app.Unsubscribe()
	defer app.Disconnect(250)

	go func() {
		for {
			select {
			case message := <-messageChannel:
				err := app.messageReceived(message)
				if err != nil {
					log.Println("error handling message: ", err)
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	log.Println("starting front end service on port 8081")
	go http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panic(err)
	}

	log.Println("starting server on port 8080")
	go srv.ListenAndServe()

	select {}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerUrl string
	}

	serverIP := os.Getenv("SERVER_IP")
	if serverIP == "" {
		serverIP = "localhost"
	}

	data.BrokerUrl = fmt.Sprintf("http://%s:%s", serverIP, "8080")

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
