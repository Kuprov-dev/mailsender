package db

import (
	"context"
	"log"
	"simplemailsender/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBUserDAO struct {
	db             *mongo.Database
	userCollection string
}

func NewMongoDBUserDAO(ctx context.Context, db *mongo.Database) *MongoDBUserDAO {
	return &MongoDBUserDAO{db: db, userCollection: "templates"}
}
func (dao *MongoDBUserDAO) PrintTemplates() {

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
		log.Println(Mail)
		log.Println(Mail.Message)
		log.Println(Mail.Status)
	}

}

//func (dao *MongoDBUserDAO) GetByUsername(ctx context.Context, username string) (*models.User, error) {
//collection := dao.db.Collection(dao.userCollection)
//
//log.Println("Searchinh by username ", username)
//// filter := bson.D{
//// 	{"comments", bson.D{{"$gt", 300}}},
//// 	{"tags", bson.D{{"$elemMatch", bson.M{"$eq": "programming"}}}},
//// }
//
//filter := bson.M{"username": username}
//var mails models.Template
//err := collection.FindOne(ctx, filter).Decode(&mails)
//
//if err != nil {
//if err == mongo.ErrNoDocuments {
//return nil, err
//}
//return nil, err
//}
//
//return &user, nil
//}

//func (dao *MongoDBUserDAO) GetByUUID(ctx context.Context, userUUID primitive.ObjectID) (*models.User, error) {
//	collection := dao.db.Collection(dao.userCollection)
//	filter := bson.M{"_id": userUUID}
//	var user models.User
//	err := collection.FindOne(ctx, filter).Decode(&user)
//
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, nil
//		}
//		return nil, err
//	}
//	return &user, nil
//}

//func (dao *MongoDBUserDAO) AddReceiver(ctx context.Context, userUUID primitive.ObjectID, recieverEmail string) error {
//	collection := dao.db.Collection(dao.userCollection)
//
//	user, err := dao.GetByUUID(ctx, userUUID)
//	if err != nil {
//		return err
//	}
//
//	fitler := bson.M{"_id": user.UUID}
//	update := bson.M{"$addToSet": bson.M{"receivers": recieverEmail}}
//	_, err = collection.UpdateOne(ctx, fitler, update)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func (dao *MongoDBUserDAO) RemoveReceiver(ctx context.Context, userUUID primitive.ObjectID, recieverEmail string) error {
//	collection := dao.db.Collection(dao.userCollection)
//
//	user, err := dao.GetByUUID(ctx, userUUID)
//	if err != nil {
//		return err
//	}
//
//	fitler := bson.M{"_id": user.UUID}
//	update := bson.M{"$pull": bson.M{"receivers": recieverEmail}}
//	_, err = collection.UpdateOne(ctx, fitler, update)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
