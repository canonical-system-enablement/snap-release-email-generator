package main

import (
	"encoding/json"
	"io/ioutil"
)

type TrelloSecrets struct {
	AppKey string `json:"app_id"`
	Token  string `json:"token"`
}

func NewTrelloSecrets(secretsFile string) (*TrelloSecrets, error) {
	trelloSecrets := new(TrelloSecrets)
	data, err := ioutil.ReadFile(secretsFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &trelloSecrets)
	if err != nil {
		return nil, err
	}

	return trelloSecrets, nil
}
