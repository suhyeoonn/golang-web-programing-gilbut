package main

import (
	"context"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/urfave/negroni/v3"
)

var (
	renderer *render.Render
	// mongoSession *mgo.Session
	mongoSession *mongo.Client
)

func init() {
	renderer = render.New()

	// s, err := mgo.Dial("127.0.0.1")
	// if err != nil {
	// 	panic(err)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	mongoSession = client
}

const (
	sessionKey    = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
)

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"host": r.Host})
	})

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "login", nil)
	})
	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// 세션에서 사용자 정보 제거 후 로그인 페이지로 이동
		sessions.GetSession(r).Delete(currentUserKey) //keyCurrentUser가 아니라 저거 아님??
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	router.GET("/auth/:action/:provider", loginHandler)

	router.POST("/rooms", createRoom)
	router.GET("/rooms", retrieveRooms)
	router.GET("/rooms/:id/messages", retrieveMessage)

	// negroni : 웹 서버의 라이프사이클을 관리하고 모든 웹 요청을 받아서 처리하는 역할
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))

	// 쿠키 기반의 세션 저장소를 만들어 negroni에서 사용할 수 있도록 등록
	// sessions.Sessions(sessionKey, store)
	n.Use(sessions.Sessions(sessionKey, store))

	n.Use(LoginRequired("/login", "/auth"))

	n.UseHandler(router)

	n.Run(":3000")
}
