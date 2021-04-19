package feedback

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createFeedback(client *mongo.Client, feedback Feedback) (interface{}, error) {

	collection := client.Database("test").Collection("feedback")
	insertResult, err := collection.InsertOne(context.TODO(), feedback)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

func getFeedbackById(client *mongo.Client, feedback_id string) (Feedback, error) {

	var feedback Feedback
	collection := client.Database("test").Collection("feedback")
	err := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: feedback_id}}).Decode(&feedback)
	if err != nil {
		log.Fatal(err)
		return Feedback{}, err
	}
	return feedback, nil
}

func getFeedbackByRootId(client *mongo.Client, root_feedback_id string) ([]Feedback, error) {
	var feedbacks []Feedback
	collection := client.Database("test").Collection("feedback")

	findOptions := options.Find()
	findOptions.SetLimit(10)

	cur, err := collection.Find(context.TODO(), bson.D{{Key: "root_feedback_id", Value: root_feedback_id}}, findOptions)
	if err != nil {
		log.Fatal(err)
		return feedbacks, err
	}

	for cur.Next(context.TODO()) {
		var feedback Feedback
		err := cur.Decode(&feedback)
		if err != nil {
			log.Fatal(err)
			return feedbacks, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return feedbacks, err
	}

	return feedbacks, nil
}

func updateFeedback(client *mongo.Client, feedback_id string, update primitive.D) (Feedback, error) {
	collection := client.Database("test").Collection("feedback")

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": feedback_id}, update)
	if err != nil {
		log.Fatal(err)
		return Feedback{}, err
	}

	return getFeedbackById(client, feedback_id)
}
