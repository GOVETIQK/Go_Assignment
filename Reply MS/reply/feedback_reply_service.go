package reply

import (
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReplyDTO struct {
	User_id          string `json:"user_id"`
	Root_feedback_id string `json:"root_feedback_id"`
	Aurthor_id       string `json:"aurthor_id"`
	Partner_id       string `json:"partner_id"`
	Comment          string `json:"comment"`
}

func ReplyFromDTO(replyDTO ReplyDTO, analysis string) Reply {

	comment := "Thank You for You Kind Words"
	if analysis == "Negative" {
		comment = "We wil try to do Better"
	}

	var reply = Reply{
		Id:               uuid.New().String(),
		Partner_id:       replyDTO.Partner_id,
		User_id:          replyDTO.User_id,
		Root_feedback_id: replyDTO.Root_feedback_id,
		Aurthor_id:       replyDTO.Aurthor_id,
		Comment:          comment,
		Analysis:         "Positve",
		Viewed:           false,
		Replied:          false,
		Last_updated_on:  primitive.NewDateTimeFromTime(time.Now()),
	}
	return reply
}

func ReplyFeedback(client *mongo.Client, replyDTO ReplyDTO, analysis string) (int, error) {
	reply := ReplyFromDTO(replyDTO, analysis)

	_, err1 := createReply(client, reply)
	if err1 != nil {
		log.Fatal(err1)
		return 0, err1
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "viewed", Value: false},
			{Key: "replied", Value: true},
			{Key: "last_updated_on", Value: primitive.NewDateTimeFromTime(time.Now())},
		},
	}}
	_, err2 := updateFeedback(client, reply.Root_feedback_id, update)
	if err2 != nil {
		log.Fatal(err2)
		return 0, err2
	}

	return 1, nil
}
