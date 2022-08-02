package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Payload struct {
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Id      int      `json:"id"`
	Clients []Client `json:"clients"`
	Files   []File   `json:"files"`
}

type Client struct {
	Id string `json:"id"`
}

type File struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	var err error
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}
	defer wsConn.Close()

	payload := readTemplate()
	json, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("could not marshal: %s\n", err.Error())
		return
	}

	err = wsConn.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		fmt.Printf("could not write: %s\n", err.Error())
		return
	}

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	router.HandleFunc("/socket", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func readTemplate() Payload {
	file, err := ioutil.ReadFile("template.json")
	if err != nil {
		log.Fatal(err)
	}
	var payload Payload
	json.Unmarshal(file, &payload)
	return payload
}
