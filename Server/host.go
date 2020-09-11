package main

import (
	"MoniTutorConnectionWatcher/Server/business"
	"MoniTutorConnectionWatcher/Server/external"
	"MoniTutorConnectionWatcher/Server/remote"
	"flag"
	"log"
	"net/http"
)

func main() {
	prefix := "/moniWatcher"

	filepath := flag.String("f", "Server/config.json", "path to config.json ")
	flag.Parse()

	user := &business.UserCred{}
	updater := &external.Updater{}

	updater, err := external.ReadConfig(updater, *filepath)

	if err != nil {
		log.Println("failed to read config File:", err.Error())
		return
	}

	sessions := make(chan business.Session, 200)

	remote.FunctionHandler(prefix, user, sessions, updater)

	go func() {
		business.Watcher(sessions, updater)
	}()

	_ = http.ListenAndServe(":8080", nil)
}