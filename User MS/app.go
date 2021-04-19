// app.go

package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/GOVETIQK/Go_Assignment/feedback"
	"github.com/GOVETIQK/Go_Assignment/util"
	"github.com/cdipaolo/sentiment"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config util.Config

type App struct {
	Router *mux.Router
	Client *mongo.Client
	Model  sentiment.Models
	ctx    context.Context
	i      int
	writer *kafka.Writer
}

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.Database.Host)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		panic("MongoDB Connection Failure")
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		panic("MongoDB Connection Failure")
	}
	return client
}

func (a *App) Initialize() {
	config := util.LoadConfiguration("app_config.json")
	a.Client = ConnectDB()
	a.Router = mux.NewRouter()
	a.Model, _ = sentiment.Restore()
	a.ctx = context.Background()
	a.i = 0
	a.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.Kafka.Broker},
		Topic:   config.Kafka.Topic,
	})
	a.InitializeRoutes()

}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":"+config.Port, a.Router))
}

func (a *App) produce(payload string) {

	err := a.writer.WriteMessages(a.ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(a.i)),
		Value: []byte(payload),
	})

	if err != nil {
		log.Fatal(err)
	}

	a.i++
}

func (a *App) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	feedback, err := feedback.SubmitFeedback(r.Body, a.Client, &a.Model)
	if err != nil {
		log.Fatal(err)
		util.RespondWithError(w, http.StatusInternalServerError, err)
	}
	a.produce(feedback.ExtarctMessage())
	util.RespondWithJSON(w, http.StatusCreated, feedback)
}

func (a *App) GetFeedbackThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	feedbacks, err := feedback.VeiwFeedbackThread(a.Client, params["feedbackId"])
	if err != nil {
		log.Fatal(err)
		util.RespondWithError(w, http.StatusInternalServerError, err)
	}
	util.RespondWithJSON(w, http.StatusOK, feedbacks)
}

func (a *App) UpdateViewedStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	feedback, err := feedback.UpdateStatus(a.Client, params["feedbackId"])
	if err != nil {
		log.Fatal(err)
		util.RespondWithError(w, http.StatusInternalServerError, err)
	}
	util.RespondWithJSON(w, http.StatusOK, feedback)
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/submitfeedback", a.SubmitFeedback).Methods("POST")
	a.Router.HandleFunc("/feedback/{feedbackId}", a.GetFeedbackThread).Methods("GET")
	a.Router.HandleFunc("/status/{feedbackId}", a.UpdateViewedStatus).Methods("PUT")
}
