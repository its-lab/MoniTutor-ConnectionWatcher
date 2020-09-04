package main

import (
	"MoniTutorConnectionWatcher/Server/business"
	"MoniTutorConnectionWatcher/Server/external"
	"MoniTutorConnectionWatcher/Server/remote"
	"net/http"
)

func main() {
	prefix := "/moniWatcher"

	user := &business.UserCred{}
	updater := &external.Updater{
		UrlPost: "http://10.0.0.10:5984/monitutor-results/_find?filter=_view&view=host_status",
		UrlBase: "http://10.0.0.10:5984/monitutor-results/",
	}

	sessions := make(chan business.Session, 200)

	remote.FunctionHandler(prefix, user, sessions, updater)

	go func() {
		business.Watcher(sessions, updater)
	}()

	_ = http.ListenAndServe(":8080", nil)
}