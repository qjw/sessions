package main

import (
	"net/http"
	"github.com/qjw/session"
	"log"
	// "github.com/gorilla/sessions"
)

var store = session.NewCookieStore([]byte("something-very-secret"))

func MyHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if aaa,ok := session.Values["foo"];ok{
		log.Print(aaa)
	}
}

func Test(w http.ResponseWriter, r *http.Request)() {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)
}

func AA(w http.ResponseWriter, r *http.Request)() {
	Test(w,r)
	MyHandler(w,r)
}

func main() {
	http.HandleFunc("/", AA)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}