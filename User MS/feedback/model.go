// Feedback Model
package feedback

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	Id               string             `json:"_id" bson:"_id"`
	Root_feedback_id string             `json:"root_feedback_id" bson:"root_feedback_id"`
	User_id          string             `json:"user_id" bson:"user_id"`
	Aurthor_id       string             `json:"aurthor_id" bson:"aurthor_id"`
	Partner_id       string             `json:"partner_id" bson:"partner_id"`
	Comment          string             `json:"comment" bson:"comment"`
	Analysis         string             `json:"analysis" bson:"analysis"`
	Viewed           bool               `json:"viewed" bson:"viewed"`
	Replied          bool               `json:"replied" bson:"replied"`
	Last_updated_on  primitive.DateTime `json:"last_updated_on" bson:"last_updated_on"`
}

func (feedback *Feedback) printFeedback() {
	fmt.Printf("%#v\n", feedback)
}

func (feedback *Feedback) ExtarctMessage() string {
	return "{ \"root_feedback_id\" :\"" + feedback.Id + "\", \"partner_id\" : \"" + feedback.Partner_id + "\", \"user_id\" : \"" +
		feedback.Aurthor_id + "\", \"analysis\" : \"" + feedback.Analysis + "\" }"
}
