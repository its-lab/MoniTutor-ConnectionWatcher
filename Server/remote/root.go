package remote

import (
	"MoniWatchDog/MoniTutorWatchDog/Server/business"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
		case http.MethodPost:
			user, err := DecodeCredentials(request)

			if err != nil {
				log.Println("failed to decode posted payload:", err.Error())
				http.Error(writer, "failed", http.StatusBadRequest)
				return
			}

			fmt.Printf("%v.%s.%v, %v:%v:%v : Creating new session with username: '%s' and hostname: '%s' \n",
				time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(),user.Username, user.Hostname)

			go business.CreateNewSess(user, sessions, updater)
		case http.MethodPut:
			user, err := DecodeCredentials(request)

			if err != nil {
				log.Println("failed to decode posted payload:", err.Error())
				http.Error(writer, "failed", http.StatusBadRequest)
				return
			}

			go business.UpdateSess(user, sessions, updater)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func DecodeCredentials(request *http.Request) (*business.UserCred, error) {
	user := &business.UserCred{}

	err := json.NewDecoder(request.Body).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}