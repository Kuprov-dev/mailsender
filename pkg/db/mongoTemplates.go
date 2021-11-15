package db

import (
	"context"
	"log"
	"simplemailsender/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBTempDAO struct {
	db             *mongo.Database
	userCollection string
}

func NewMongoDBTemp(ctx context.Context, db *mongo.Database) *MongoDBTempDAO {
	return &MongoDBTempDAO{db: db, userCollection: "templates"}
}
func (dao *MongoDBTempDAO) PrintTemplates() {

	collection := dao.db.Collection("templates")

	var mails []models.Template
	var mail models.Template
	cur, err := collection.Find(context.TODO(), bson.M{})

	for cur.Next(context.TODO()) {
		err = cur.Decode(&mail)
		if err != nil {

			panic(err)
			return
		}

		if mail.Status == "pending" {
			mails = append(mails, mail)
		}

	}
	for _, Mail := range mails {
		Mail.Status = "done"
		_, err = collection.UpdateOne(context.TODO(), bson.D{{"status", bson.D{{"$eq", "pending"}}}}, bson.D{{"$set", bson.D{{"status", "done"}}}})
		if err != nil {
			panic(err)

		}

	}
	for _, Mail := range mails {
		Mail.Status = "done"
		log.Println(Mail)
		log.Println(Mail.Message)

		log.Println(Mail.Status)
	}

}
