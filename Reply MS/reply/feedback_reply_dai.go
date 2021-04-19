package reply

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createReply(client *mongo.Client, reply Reply) (bool, error) {

	collection := client.Database("test").Collection("feedback")
	_, err := collection.InsertOne(context.TODO(), reply)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return true, nil
}

func updateFeedback(client *mongo.Client, reply_id string, update primitive.D) (bool, error) {

	collection := client.Database("test").Collection("feedback")
	updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": reply_id}, update)

	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return updateResult.MatchedCount == 1, nil
}
