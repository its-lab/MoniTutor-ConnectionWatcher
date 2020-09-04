package remote

import (
	"MoniTutorConnectionWatcher/Server/business"
	"net/http"
)

var usercred *business.UserCred
var sessions chan business.Session
var updater business.CouchConnector

func FunctionHandler(prefix string, _userSess *business.UserCred, _sessionsOrig chan business.Session, _updater business.CouchConnector) {
	usercred = _userSess
	sessions = _sessionsOrig
	updater = _updater

	http.HandleFunc(prefix+"/", rootHandler)
}