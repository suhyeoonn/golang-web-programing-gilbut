package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

const messageFetchSize = 10

type Message struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	RoomId    bson.ObjectId `bson:"room_id" json:"room_id"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	User      *User         `bson:"user" json:"user"`
}

func (m *Message) create() error {
	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()

	c := mongoSession.Database("test").Collection("messages")

	_, err := c.InsertOne(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}

func retrieveMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("retrieveMessage")
}
