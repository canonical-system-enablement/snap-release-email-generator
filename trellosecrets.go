/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
