package db

import (
	"context"
	"log"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBTemplateDAO struct {
	db                 *mongo.Database
	templateCollection string
	opt                *options.ClientOptions
}

var mongoDBConnection *mongo.Database
var optionsDB *options.ClientOptions

func NewMongoDBTemplateDAO(ctx context.Context, db *mongo.Database) *MongoDBTemplateDAO {
	return &MongoDBTemplateDAO{db: db, templateCollection: "templates", opt: optionsDB}
}
func ConnectMongoDB(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	optionsDB = options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoDBConnection = client.Database("mydb")
	log.Println("Connection to DB success")
}

func GetMongoDBConnection() *mongo.Database {
	return mongoDBConnection
}
