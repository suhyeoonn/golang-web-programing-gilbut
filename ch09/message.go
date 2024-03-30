package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

	fmt.Printf("왜 안돼 %+v", m)
	c := mongoSession.Database("test").Collection("messages")

	result, err := c.InsertOne(context.Background(), m)
	fmt.Println("why........", result)
	if err != nil {
		fmt.Println("ERR!!!!!", err)
		return err
	}
	return nil
}

func retrieveMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 쿼리 파라미터로 전달된 limit 값 확인
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	fmt.Println("limit: ", limit)
	if err != nil {
		// 정상적인 limit 값이 전달되지 않은 경우 limit를 messageFetchSize으로 셋팅
		limit = messageFetchSize
	}

	var messages []Message
	// TODO messages 에 담기
	// _id 역순으로 정렬하여 limit 수만큼 message 조회
	c := mongoSession.Database("test").Collection("messages")
	cur, err := c.Find(context.TODO(), bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))})
	//.Sort("-_id").Limit(limit).All(&messages)
	// err = c.Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).
	// 	Sort("-_id").Limit(limit).All(&messages)
	if err != nil {
		// 오류 발생시 500 에러 리턴
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	if err = cur.All(context.TODO(), &messages); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	// message 조회 결과 리턴
	renderer.JSON(w, http.StatusOK, messages)
}
