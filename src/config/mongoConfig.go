package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

type AppConfig struct {
	Env             string
	MongoClientUri  string
	DBName          string
	UserColName     string
	TicketColName   string
	CategoryColName string
	MongoDuration   int16
	MaxPageLimit    int
}

var EnvConfig = map[string]AppConfig{
	"qa": {
		Env:             "qa",
		MongoClientUri:  "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority",
		DBName:          "TicketApp",
		UserColName:     "USer",
		TicketColName:   "Ticket",
		CategoryColName: "Category",
		MongoDuration:   5,
		MaxPageLimit:    100,
	},
	"prod": {
		Env:             "qa",
		MongoClientUri:  "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority",
		DBName:          "TicketApp",
		UserColName:     "USer",
		TicketColName:   "Ticket",
		CategoryColName: "Category",
		MongoDuration:   5,
		MaxPageLimit:    100,
	},
}

func NewMongoConfig() AppConfig {
	return AppConfig{}
}

func (mCfg *AppConfig) CloseConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (mCfg *AppConfig) ConnectDatabase() (*mongo.Client, context.Context, context.CancelFunc, *AppConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	cfg := GetConfigModel()
	//conStr := "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority"
	if cfg.MongoClientUri == "" {
		panic("Connection string was not found. Check the .env file.")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoClientUri))
	if err != nil {
		panic(err)
	}

	return client, ctx, cancel, &cfg
}

func (mCfg *AppConfig) Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected successfully.")
	return nil
}

func (mCfg *AppConfig) GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("TicketApp").Collection(collectionName)
}

func GetConfigModel() AppConfig {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(wd + "/src/.env")

	if err != nil {
		panic("Env was not found.")
		log.Fatalln("error .env")
	}

	env := os.Getenv("Env")
	if env == "" {
		panic("Env was not found.")
	}

	model, exist := EnvConfig[env]
	if !exist {
		panic("There is no model with provided environment.")
	}

	return model
}
