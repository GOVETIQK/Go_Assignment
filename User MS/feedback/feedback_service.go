package feedback

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/cdipaolo/sentiment"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackDTO struct {
	User_id          string `json:"user_id"`
	Root_feedback_id string `json:"root_feedback_id"`
	Aurthor_id       string `json:"aurthor_id"`
	Partner_id       string `json:"partner_id"`
	Comment          string `json:"comment"`
}

func FeedbackFromDTO(feedbackDTO FeedbackDTO, model *sentiment.Models) Feedback {

	var feedback = Feedback{
		Id:               uuid.New().String(),
		Partner_id:       feedbackDTO.Partner_id,
		User_id:          feedbackDTO.User_id,
		Root_feedback_id: feedbackDTO.Root_feedback_id,
		Aurthor_id:       feedbackDTO.Aurthor_id,
		Comment:          feedbackDTO.Comment,
		Analysis:         Analyze(feedbackDTO.Comment, model),
		Viewed:           false,
		Replied:          false,
		Last_updated_on:  primitive.NewDateTimeFromTime(time.Now()),
	}
	return feedback
}

func Analyze(content string, model *sentiment.Models) string {
	if model.SentimentAnalysis(content, sentiment.English).Score == 1 {
		return "Positive"
	} else {
		return "Negative"
	}
}

func SubmitFeedback(request_body io.ReadCloser, client *mongo.Client, model *sentiment.Models) (Feedback, error) {
	var feedbackDTO FeedbackDTO

	if err := json.NewDecoder(request_body).Decode(&feedbackDTO); err != nil {
		log.Fatal(err)
		return Feedback{}, err
	}

	feedback := FeedbackFromDTO(feedbackDTO, model)

	_, err := createFeedback(client, feedback)
	if err != nil {
		log.Fatal(err)
		return Feedback{}, err
	}

	return feedback, nil
}

func VeiwFeedbackThread(client *mongo.Client, feedback_id string) ([]Feedback, error) {
	var feedbacks []Feedback

	headPost, err := getFeedbackById(client, feedback_id)
	if err != nil {
		log.Fatal(err)
		return feedbacks, err
	}
	feedbacks = append(feedbacks, headPost)

	replies, err := getFeedbackByRootId(client, feedback_id)
	if err != nil {
		log.Fatal(err)
		return feedbacks, err
	}
	feedbacks = append(feedbacks, replies...)

	return feedbacks, nil
}

func UpdateStatus(client *mongo.Client, feedback_id string) (Feedback, error) {
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "viewed", Value: true},
		},
	}}
	return updateFeedback(client, feedback_id, update)
}
