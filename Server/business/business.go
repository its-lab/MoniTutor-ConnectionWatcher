package business

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

var ErrConnectionLoss = errors.New("lost connection")
var ErrUnexpectedStatusCode = errors.New("unexpected Status Code")

type Docs struct {
	Id string `json:"_id"`
	Rev string `json:"_rev"`
	Username string `json:"username"`
	SeverityCode int `json:"severity_code"`
	Hostname string `json:"hostname"`
	Time string `json:"time"`
	Output string `json:"output"`
	Type string `json:"type"`
	Address string `json:"address"`
}

type Result struct {
	Documents[] Docs `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning string `json:"warning"`
}

type Params struct {
	Selector Selector `json:"selector"`
}

type Selector struct {
	Id ID `json:"_id"`
}

type ID struct {
	Eq string `json:"$eq"`
}

type UserCred struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}

type Session struct {
	User UserCred
	timeStamp time.Time
	checked bool
}

type CouchConnector interface {
	GetRev(body io.Reader) (*Result, error)
	UpdateDocument(docs *Result, id string, output string, severityCode int) error
}

func buildBody(params *Params) (io.Reader, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(params)

	if err != nil {
		log.Println("failed to build body from data struct:", err.Error())
		return nil, err
	}

	return buffer, nil
}

func CreateNewSess(user *UserCred, sessions chan Session, updater CouchConnector) {
	session := Session {
		User: *user,
		timeStamp: time.Now(),
		checked: false,
	}

	id := user.Username + "_"  + user.Hostname

	params := &Params{
		Selector: Selector{
			Id: ID{Eq: id},
		},
	}

	body, err := buildBody(params)

	if err != nil {
		log.Println("failed to build request body: ", err.Error())
	}

	docs, err := updater.GetRev(body)

	if err != nil {
		log.Println("failed to get rev number")
	}

	err = updater.UpdateDocument(docs, id, "Connected", 0)

	if err != nil {
		log.Println("failed to update document")
	}

	sessions<-session
}

func UpdateSess(user *UserCred, sessions chan Session, updater CouchConnector) {
	for sess := range sessions {
		if strings.EqualFold(sess.User.Hostname, user.Hostname) && strings.EqualFold(sess.User.Username, user.Username) {
			CreateNewSess(&UserCred{sess.User.Username, sess.User.Hostname}, sessions, updater)
			return
		}
	}
}

func Watcher(sessions chan Session, updater CouchConnector) {
	for range time.Tick(30 * time.Second) {
		for sess := range sessions {
			runtime.Gosched()

			if sess.checked {
				sessions<-sess
				break
			}

			err := Watch(sess, sessions)

			if err != nil {
				id := sess.User.Username + "_"  + sess.User.Hostname

				params := &Params{
					Selector: Selector{
						Id: ID{Eq: id},
					},
				}

				body, err := buildBody(params)

				if err != nil {
					log.Println("failed to build request body: ", err.Error())
				}

				fmt.Println(body)

				docs, err := updater.GetRev(body)

				if err != nil {
					log.Println("failed to get rev number")
				}

				err = updater.UpdateDocument(docs, id, "Disconnected", 1)

				if err != nil {
					log.Println("failed to update document")
				}

				//err = updater.ReloadMoniTutor(sess.User.Username, "%")
			}
		}

		for sess := range sessions {
			runtime.Gosched()
			if !sess.checked {
				sessions<-sess
				break
			}

			sess.checked = false
			sessions<-sess
		}
	}
}

func Watch(sess Session, sessions chan Session) error {
	diff := time.Now().Sub(sess.timeStamp)

	if int(diff.Seconds()) > 120 {
		fmt.Printf("%v.%s.%v, %v:%v:%v : Lost connection to machine with username = '%s' and hostname = '%s' \n",
			time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), sess.User.Username, sess.User.Hostname)
		return ErrConnectionLoss
	} else if int(diff.Seconds()) > 30 {
		fmt.Printf("%v.%s.%v, %v:%v:%v : Searching for connection to machine with username = '%s' and hostname = '%s': timeStamp = '%d' \n",
			time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), sess.User.Username, sess.User.Hostname, int(diff.Seconds()))
		sess.checked = true
		sessions<-sess
		return nil
	}

	fmt.Printf("%v.%s.%v, %v:%v:%v : Stable Connection to machine with username = '%s' and hostname = '%s': timeStamp = '%d' \n",
		time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), sess.User.Username, sess.User.Hostname, int(diff.Seconds()))

	sess.checked = true
	sessions<-sess

	return nil
}

