package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni/v3"
)

const (
	nextPageKey     = "next_page"
	authSecurityKey = "auth_security_key"
)

func init() {
	fmt.Println("auth init!!")
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New(clientID, clientSecret, "http://localhost:3000/auth/callback/google"),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	// s := sessions.GetSession(r)

	switch action {
	case "login":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("login", provider, loginUrl)
		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		fmt.Println("-==================================")
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}
		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("callback", user)

		u := &User{
			Uid:      user.Data().Get("id").MustStr(),
			Name:     user.Name(),
			Email:    user.Email(),
			AvataUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u)
		// http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
		http.Redirect(w, r, "/login", http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"' is not supported", http.StatusNotFound)
	}
}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				fmt.Println("ignore~~~~")
				next(rw, r)
				return
			}
		}
		u := GetCurrentUser(r)

		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(rw, r)
			return
		}

		SetCurrentUser(r, nil)
		// 현재 URL을 세션에 저장하고 로그인 페이지로 리다이렉트. 로그인을 성공하면 세션에 저장한 URL로 다시 리다이렉트
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())
		http.Redirect(rw, r, "/login", http.StatusFound)
	}
}
