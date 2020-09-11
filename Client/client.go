
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Id struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}

func main() {
	user := flag.String("u", "", "MoniTutor username")
	host := flag.String("h", "", "Hostname of the system (itsclient/itsserver/itsjumphost)")
	ip := flag.String("a", "", "IP address of the MoniTutor System")

	flag.Parse()

	if strings.EqualFold(*user, " ") || strings.EqualFold(*host, " ") {
		log.Printf("invalid username or hostname. Username: '%s', Hostname: '%s'", *user, *host)
		return
	}

	id := getId(user, host)

	err := run(id, ip)

	if err != nil {
		log.Println("failed to run MoniWatchDog:", err.Error())
	}
}

func run(id *Id, ip *string) error {
	_, err := authenticate(id, ip)

	if err != nil {
		log.Println("failed to authenticate to MoniWatchDog:", err.Error())
		return err
	}

	counter := 0

	for range time.Tick(30 * time.Second) {
		resp, isTimeout, err := sendWoof(id, *ip)

		if err != nil && !isTimeout {
			log.Println("failed to send request:", err.Error())

			return err
		}

		if isTimeout {
			counter++
		}

		if resp != nil {
			if resp.StatusCode != 200 {
				log.Println("failed to send woof to server, Status Code:", resp.StatusCode)
				counter++
			} else if resp.StatusCode == 200 {
				counter = 0
			}
		}

		for counter > 4 {
			log.Println("lost connection to server, trying to reconnect...")

			isReconnected, err := authenticate(id, ip)

			if err != nil {
				log.Println("failed to reconnect")
			}

			if isReconnected {
				counter = 0
				break
			}

			time.Sleep(10 * time.Second)
		}
	}

	return nil
}

func getId(user *string, host *string) *Id {
	id := &Id{}

	id.Username = *user
	id.Hostname = *host

	return id
}

func authenticate(id *Id, ip *string) (bool, error) {
	body, err := buildBody(id)

	if err != nil {
		log.Println("failed to build body for authentication post:", err.Error())
		return false, err
	}

	resp, err := http.Post("http://" + *ip + ":8080/moniWatcher/", "application/json", body)

	if err != nil {
		log.Println("failed to send authentication post:", err.Error())
		return true, err
	}

	fmt.Println(resp.Body)

	return true, nil
}

func sendWoof(id *Id, ip string) (*http.Response, bool, error) {
	body, err := buildBody(id)

	if err != nil {
		log.Println("failed to get message body:", err.Error())
		return nil, false, err
	}

	request, err := http.NewRequest("PUT", "http://" + ip + ":8080/moniWatcher/sessions/", body)

	if err != nil {
		log.Println("failed to build http request:", err.Error())
		return nil, false, err
	}

	fmt.Println(request.Body)

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(request)

	if err, ok := err.(net.Error); ok {
		log.Println("failed to send request, missing network connection:", err.Timeout())

		return resp, true, nil
	} else if err, ok := err.(net.Error); !ok {
		//log.Println("failed to send request, missing network connection:", err.Timeout())

		return resp, true, nil
	} else if err != nil {
		log.Println("failed to send request:", err.Error())

		return resp, false, err
	} else {
		return resp, false, nil
	}
}

func buildBody(id *Id) (io.Reader, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(id)

	if err != nil {
		log.Println("failed to build json body from data struct:", err.Error())
		return nil, err
	}

	return buffer, nil
}
