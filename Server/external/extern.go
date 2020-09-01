package external

import (
	"MoniWatchDog/MoniTutorWatchDog/Server/business"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Updater struct {
	UrlPost string
	UrlBase string
}



func (m *Updater) GetRev(body io.Reader) (*business.Result, error) {
	result := &business.Result{}

	request, err := http.NewRequest("POST", m.UrlPost, body)

	if err != nil {
		log.Println("failed to build http request:", err.Error())
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	timeout := time.Duration(7 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Println("failed to send request:", err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("Status Code != 200. Status Code: ", resp.StatusCode)

		return nil, business.ErrUnexpectedStatusCode
	}

	err = json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		fmt.Println("failed to decode payload:", err.Error())
		return nil, err
	}

	fmt.Println(result)

	return result, nil
}

func (m *Updater) UpdateDocument(docs *business.Result, id string, output string, severityCode int) error {
	rev := docs.Documents[0].Rev

	docs.Documents[0].Output = output
	docs.Documents[0].SeverityCode = severityCode

	buffer := new(bytes.Buffer)

	err := json.NewEncoder(buffer).Encode(docs.Documents[0])

	if err != nil {
		log.Println("failed to encode payload:", err.Error())
		return err
	}

	request, err := http.NewRequest("PUT", m.UrlBase + id, buffer)

	if err != nil {
		log.Println("failed to create update request:", err.Error())
		return err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("If-Match", rev)

	timeout := time.Duration(7 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Println("failed to send update request:", err)
		return err
	}

	fmt.Println(resp)

	return nil
}

func (m *Updater) ReloadMoniTutor(username string, scenarioNum string) error {
	//request, err := http.NewRequest("GET", "https://10.0.0.10/MoniTutor/scenarios/progress.html/" + scenarioNum + "/" + username, nil)

	return nil
}
