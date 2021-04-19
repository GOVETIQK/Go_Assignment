// app.go

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GOVETIQK/Go_Assignment/reply"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	topic         = "feedback"
	brokerAddress = "localhost:9092"
)

type App struct {
	Router *mux.Router
	Client *mongo.Client
	ctx    context.Context
	reader *kafka.Reader
}

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (a *App) Initialize() {
	a.Client = ConnectDB()
	a.Router = mux.NewRouter()
	a.ctx = context.Background()
	a.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	a.InitializeRoutes()

}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type Msg struct {
	Root_feedback_id string `json:"root_feedback_id"`
	Partner_id       string `json:"partner_id"`
	User_id          string `json:"user_id"`
	Analysis         string `json:"analysis"`
}

func (a *App) Consume() (reply.ReplyDTO, string) {
	var temp Msg

	msg, err := a.reader.ReadMessage(a.ctx)
	if err != nil {
		panic("could not read message " + err.Error())
	}

	err = json.Unmarshal(msg.Value, &temp)
	if err != nil {
		log.Fatal(err)
	}

	var replyDTO = reply.ReplyDTO{
		User_id:          temp.User_id,
		Root_feedback_id: temp.Root_feedback_id,
		Partner_id:       temp.Partner_id,
		Aurthor_id:       "",
		Comment:          "",
	}

	return replyDTO, temp.Analysis
}

func (a *App) ReplyToFeedbacks(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var replies []reply.Reply

	// count, err := strconv.Atoi(params["count"])
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for i := 0; i < 1; i++ {
		replyDTO, analysis := a.Consume()

		reply.ReplyFeedback(a.Client, replyDTO, analysis)
	}

	respondWithJSON(w, http.StatusCreated, replies)
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/autoReply", a.ReplyToFeedbacks).Methods("PUT")

}
