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
	"flag"
	"fmt"
	"log"
	"strings"
)

var (
	trelloSecretsFile = flag.String("secrets", "trello_secrets.json", "Trello Secrets configuration")
	snapPublisher     = flag.String("publisher", "Simon", "Snap publisher's name")
)

// Card title has a format of "NAME - VERSION"
func getSnapNameAndVersion(Name string) (string, string) {
	snapName := ""
	snapVersion := "null"

	tmp := strings.Split(Name, " - ")

	if len(tmp) > 0 {
		snapName = tmp[0]
	}

	if len(tmp) > 1 {
		snapVersion = tmp[1]
	}

	return snapName, snapVersion
}

func main() {
	flag.Parse()

	trelloSecrets, err := NewTrelloSecrets(*trelloSecretsFile)
	if err != nil {
		log.Fatal(err)
	}

	trelloProxy := TrelloProxy{}
	trelloProxy.Connect(*trelloSecrets)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := trelloProxy.CardsOfSnapsToRelease()
	if err != nil {
		log.Fatal(err)
	}

	msgBuilder := NewMessageBuilder()
	msgBuilder.SetPublisher(*snapPublisher)

	for _, c := range cards {
		// Get the data from card
		snapName, snapVersion := getSnapNameAndVersion(c.Name)

		checklists, err := c.Checklists()
		if err != nil {
			log.Fatal(err)
		}

		for _, clist := range checklists {

			if strings.Contains(clist.Name, "Changes") {
				for _, citem := range clist.CheckItems {
					msgBuilder.AddChange(snapName, snapVersion, citem.Name)
				}
			}

			if strings.Contains(clist.Name, "Bileto") {
				for _, bitem := range clist.CheckItems {
					msgBuilder.AddBiletoURL(snapName, bitem.Name)
				}
			}

		}
	}

	email := msgBuilder.ConstructMessage()
	fmt.Println(email)
}
