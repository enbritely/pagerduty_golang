package pagerduty_golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Endpoint = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"
)

type Pagerduty struct {
	Servicekey, Client, ClientUrl string
}

func responseError(resp *http.Response) (err error) {
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("%s: %s", resp.Status, string(body))
}

func (p *Pagerduty) CreateIncident(key, description string, prio int) (incidentKey string, err error) {
	return p.CreateIncidentWithDetails(key, description, prio, map[string]interface{}{})
}

func (p *Pagerduty) CreateIncidentWithDetails(key, description string, prio int, details map[string]interface{}) (incidentKey string, err error) {
	details["prio"] = prio
	payload, err := p.createJson(p.Servicekey, description, key, details)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(Endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Println("[Pagerduty] Failed to execute request")
		return "", err
	}
	defer resp.Body.Close()

	err = responseError(resp)
	if err != nil {
		log.Println("[Pagerduty] Response was an Error response")
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Pagerduty] Failed to read request body")
		return "", err
	}

	respBody := map[string]string{}
	err = json.Unmarshal(bodyBytes, &respBody)
	if err != nil {
		return "", err
	}
	return respBody["incident_key"], nil
}

func (p *Pagerduty) createJson(servicekey, description, incidentKey string, details map[string]interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"service_key": p.Servicekey,
		"event_type":  "trigger",
		"description": description,
		"details":     details,
		"client":      p.Client,
		"client_url":  p.ClientUrl,
	}
	if len(incidentKey) > 0 {
		payload["incident_key"] = incidentKey
	}
	return json.Marshal(payload)
}
