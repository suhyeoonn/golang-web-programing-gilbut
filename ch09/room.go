package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Println("create Room")
	r := new(Room)
	errs := binding.Bind(req, r)
	if errs.Handle(w) {
		return
	}

	r.ID = bson.NewObjectId()
	c := mongoSession.Database("test").Collection("rooms")

	fmt.Println(r)
	_, err := c.InsertOne(context.Background(), r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusCreated, r)
}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Println("retrieve Rooms")

	var rooms []Room
	c := mongoSession.Database("test").Collection("rooms")
	cur, err := c.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("=====error")
		log.Fatal(err)
	}
	// defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &rooms); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, rooms)
}
