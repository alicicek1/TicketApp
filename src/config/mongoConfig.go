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

type MongoConfig struct {
}

func NewMongoConfig() MongoConfig {
	return MongoConfig{}
}

func (mCfg *MongoConfig) CloseConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (mCfg *MongoConfig) ConnectDatabase() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	//conStr := GetConnectionString()
	conStr := "mongodb+srv://admin:1@cluster0.ymrmq.mongodb.net/?retryWrites=true&w=majority"
	if conStr == "" {
		panic("Connection string was not found. Check the .env file.")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conStr))
	if err != nil {
		panic(err)
	}

	return client, ctx, cancel
}

func (mCfg *MongoConfig) Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("Connected successfully.")
	return nil
}

func (mCfg *MongoConfig) GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("TicketApp").Collection(collectionName)
}

func GetConnectionString() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("error .env")
	}

	mongoIRU := os.Getenv("ConnectionString")
	return mongoIRU
}
